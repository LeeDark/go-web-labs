package data

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
)

func TestMovieModelInsertAndGetWithPostgres(t *testing.T) {
	dsn := os.Getenv("GREENLIGHT_TEST_DB_DSN")
	if dsn == "" {
		t.Skip("GREENLIGHT_TEST_DB_DSN is not set; skipping PostgreSQL integration test")
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

func applyMovieMigrations(t *testing.T, db *sql.DB) {
	t.Helper()
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
		if _, err := db.Exec(string(sqlBytes)); err != nil {
			t.Fatalf("apply migration %s: %v", name, err)
		}
	}
}
