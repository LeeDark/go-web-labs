package data

import (
	"reflect"
	"strings"
	"testing"

	"github.com/LeeDark/go-web-labs/books/lets-go-further/greenlight/internal/validator"
)

func TestValidateMovie(t *testing.T) {
	tests := []struct {
		name       string
		movie      Movie
		wantErrors map[string]string
	}{
		{
			name: "valid movie",
			movie: Movie{
				Title:   "The Go Programming Language",
				Year:    2015,
				Runtime: 380,
				Genres:  []string{"technical", "programming"},
			},
			wantErrors: map[string]string{},
		},
		{
			name:  "missing required values",
			movie: Movie{},
			wantErrors: map[string]string{
				"title":   "must be provided",
				"year":    "must be provided",
				"runtime": "must be provided",
				"genres":  "must be provided",
			},
		},
		{
			name: "year before the supported range",
			movie: Movie{
				Title:   "Early cinema",
				Year:    1888,
				Runtime: 1,
				Genres:  []string{"history"},
			},
			wantErrors: map[string]string{"year": "must be greater than 1888"},
		},
		{
			name: "minimum supported year",
			movie: Movie{
				Title:   "Early cinema",
				Year:    1889,
				Runtime: 1,
				Genres:  []string{"history"},
			},
			wantErrors: map[string]string{},
		},
		{
			name: "future year",
			movie: Movie{
				Title:   "Future release",
				Year:    3000,
				Runtime: 90,
				Genres:  []string{"science fiction"},
			},
			wantErrors: map[string]string{"year": "must not be in the future"},
		},
		{
			name: "title exceeds limit",
			movie: Movie{
				Title:   strings.Repeat("a", 501),
				Year:    2000,
				Runtime: 90,
				Genres:  []string{"drama"},
			},
			wantErrors: map[string]string{"title": "must not be more than 500 bytes long"},
		},
		{
			name: "negative runtime",
			movie: Movie{
				Title:   "Impossible duration",
				Year:    2000,
				Runtime: -1,
				Genres:  []string{"drama"},
			},
			wantErrors: map[string]string{"runtime": "must be a positive integer"},
		},
		{
			name: "duplicate genres",
			movie: Movie{
				Title:   "Repeated category",
				Year:    2000,
				Runtime: 90,
				Genres:  []string{"drama", "drama"},
			},
			wantErrors: map[string]string{"genres": "must not contain duplicate value"},
		},
		{
			name: "empty genres",
			movie: Movie{
				Title:   "No genre",
				Year:    2000,
				Runtime: 90,
				Genres:  []string{},
			},
			wantErrors: map[string]string{"genres": "must contain at least 1 genre"},
		},
		{
			name: "too many genres",
			movie: Movie{
				Title:   "Every genre",
				Year:    2000,
				Runtime: 90,
				Genres:  []string{"action", "comedy", "drama", "fantasy", "horror", "thriller"},
			},
			wantErrors: map[string]string{"genres": "must not contain more than 5 genres"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateMovie(v, &tt.movie)

			if !reflect.DeepEqual(v.Errors, tt.wantErrors) {
				t.Fatalf("validation errors = %#v, want %#v", v.Errors, tt.wantErrors)
			}
		})
	}
}
