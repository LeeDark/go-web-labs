package data

import (
	"database/sql"
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
	defer db.Close()

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

	if err := model.Insert(movie); err != nil {
		t.Fatalf("insert movie: %v", err)
	}
	if movie.ID < 1 || movie.CreatedAt.IsZero() || movie.Version != 1 {
		t.Fatalf("generated fields = id %d, created_at %v, version %d", movie.ID, movie.CreatedAt, movie.Version)
	}

	defer func() {
		if _, err := db.Exec("DELETE FROM movies WHERE id = $1", movie.ID); err != nil {
			t.Errorf("cleanup movie %d: %v", movie.ID, err)
		}
	}()

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
}

func applyMovieMigrations(t *testing.T, db *sql.DB) {
	t.Helper()
	for _, name := range []string{
		"000001_create_movie_table.up.sql",
		"000002_add_movies_check_constraints.up.sql",
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
