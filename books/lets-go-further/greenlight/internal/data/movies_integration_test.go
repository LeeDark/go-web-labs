package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

func TestMovieModelInsertAndGetWithPostgres(t *testing.T) {
	dsn := os.Getenv("GREENLIGHT_TEST_DB_DSN")
	if dsn == "" {
		t.Skip("GREENLIGHT_TEST_DB_DSN is not set; skipping PostgreSQL integration test")
	}
	expectedDatabase := os.Getenv("GREENLIGHT_TEST_DB_NAME")
	if expectedDatabase == "" {
		t.Fatal("GREENLIGHT_TEST_DB_NAME is not set; refusing to use an unnamed test database")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close test database: %v", err)
		}
	})

	if err := db.Ping(); err != nil {
		t.Fatalf("ping test database: %v", err)
	}
	verifyTestDatabase(t, db, expectedDatabase)
	applyMovieMigrations(t, db)

	movie := &Movie{
		Title:   "Stage 5 integration movie",
		Year:    2020,
		Runtime: 95,
		Genres:  []string{"drama", "testing"},
	}
	model := MovieModel{DB: db}
	invalidYear := &Movie{Title: "Invalid boundary", Year: 1888, Runtime: 90, Genres: []string{"testing"}}
	if err := model.Insert(invalidYear); err == nil {
		if _, cleanupErr := db.Exec("DELETE FROM movies WHERE id = $1", invalidYear.ID); cleanupErr != nil {
			t.Errorf("cleanup unexpectedly inserted movie %d: %v", invalidYear.ID, cleanupErr)
		}
		t.Fatal("insert movie with year 1888 succeeded; want database constraint error")
	}

	if err := model.Insert(movie); err != nil {
		t.Fatalf("insert movie: %v", err)
	}
	t.Cleanup(func() {
		if _, err := db.Exec("DELETE FROM movies WHERE id = $1", movie.ID); err != nil {
			t.Errorf("cleanup movie %d: %v", movie.ID, err)
		}
	})
	if movie.ID < 1 || movie.CreatedAt.IsZero() || movie.Version != 1 {
		t.Fatalf("generated fields = id %d, created_at %v, version %d", movie.ID, movie.CreatedAt, movie.Version)
	}

	got, err := model.Get(movie.ID)
	if err != nil {
		t.Fatalf("get movie: %v", err)
	}
	if got.ID != movie.ID || got.Title != movie.Title || got.Runtime != movie.Runtime || got.Version != movie.Version {
		t.Fatalf("movie = %+v, want fields from %+v", got, movie)
	}
	if !got.CreatedAt.Equal(movie.CreatedAt) {
		t.Fatalf("created_at = %v, want %v", got.CreatedAt, movie.CreatedAt)
	}
	if len(got.Genres) != len(movie.Genres) || got.Genres[0] != movie.Genres[0] || got.Genres[1] != movie.Genres[1] {
		t.Fatalf("genres = %#v, want %#v", got.Genres, movie.Genres)
	}

	stale, err := model.Get(movie.ID)
	if err != nil {
		t.Fatalf("get stale movie copy: %v", err)
	}
	got.Title = "Stage 5 updated movie"
	if err := model.Update(got); err != nil {
		t.Fatalf("update current movie: %v", err)
	}
	if got.Version != movie.Version+1 {
		t.Fatalf("updated version = %d, want %d", got.Version, movie.Version+1)
	}
	if err := model.Update(stale); !errors.Is(err, ErrEditConflict) {
		t.Fatalf("stale update error = %v, want %v", err, ErrEditConflict)
	}
}

func verifyTestDatabase(t *testing.T, db *sql.DB, expectedDatabase string) {
	t.Helper()

	var databaseName, currentUser, owner string
	err := db.QueryRow(`
		SELECT current_database(), current_user, pg_get_userbyid(datdba)
		FROM pg_database
		WHERE datname = current_database()`).Scan(&databaseName, &currentUser, &owner)
	if err != nil {
		t.Fatalf("inspect test database identity: %v", err)
	}
	if databaseName != expectedDatabase {
		t.Fatalf("connected database = %q, want GREENLIGHT_TEST_DB_NAME %q", databaseName, expectedDatabase)
	}
	if !strings.HasSuffix(databaseName, "_test") {
		t.Fatalf("test database name %q must end with _test", databaseName)
	}
	t.Logf("using test database %q (owner=%q, connected user=%q)", databaseName, owner, currentUser)
}

func applyMovieMigrations(t *testing.T, db *sql.DB) {
	t.Helper()

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS stage5_test_migrations (
			name text PRIMARY KEY,
			checksum text NOT NULL,
			applied_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
		)`); err != nil {
		t.Fatalf("create test migration ledger: %v", err)
	}

	for _, name := range []string{
		"000001_create_movie_table.up.sql",
		"000002_add_movies_check_constraints.up.sql",
		"000003_align_movie_year_constraint.up.sql",
	} {
		path := filepath.Join("..", "..", "migrations", name)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", path, err)
		}

		checksum := fmt.Sprintf("%x", sha256.Sum256(sqlBytes))
		var appliedChecksum string
		err = db.QueryRowContext(ctx, `SELECT checksum FROM stage5_test_migrations WHERE name = $1`, name).Scan(&appliedChecksum)
		switch {
		case errors.Is(err, sql.ErrNoRows):
			applyMigration(t, db, name, checksum, string(sqlBytes))
		case err != nil:
			t.Fatalf("read migration ledger for %s: %v", name, err)
		case appliedChecksum != checksum:
			t.Fatalf("migration %s checksum changed after it was applied; use a fresh disposable database", name)
		}
	}
}

func applyMigration(t *testing.T, db *sql.DB, name, checksum, statement string) {
	t.Helper()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		t.Fatalf("begin migration %s: %v", name, err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Errorf("rollback migration %s: %v", name, err)
		}
	}()

	if _, err := tx.Exec(statement); err != nil {
		t.Fatalf("apply migration %s: %v", name, err)
	}
	if _, err := tx.Exec(`INSERT INTO stage5_test_migrations (name, checksum) VALUES ($1, $2)`, name, checksum); err != nil {
		t.Fatalf("record migration %s: %v", name, err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit migration %s: %v", name, err)
	}
}
