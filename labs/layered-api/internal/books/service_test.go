package books

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	books     []Book
	getBook   Book
	getErr    error
	listErr   error
	createErr error
	created   CreateBookInput
}

func (r *fakeRepository) List(context.Context) ([]Book, error)         { return r.books, r.listErr }
func (r *fakeRepository) GetByID(context.Context, int64) (Book, error) { return r.getBook, r.getErr }
func (r *fakeRepository) Create(_ context.Context, input CreateBookInput) (Book, error) {
	r.created = input
	if r.createErr != nil {
		return Book{}, r.createErr
	}
	return Book{ID: 3, Title: input.Title, Author: input.Author}, nil
}

func TestServiceCreate(t *testing.T) {
	repositoryErr := errors.New("database unavailable")
	tests := []struct {
		name    string
		repo    *fakeRepository
		input   CreateBookInput
		wantErr error
	}{
		{name: "creates trimmed book", repo: &fakeRepository{}, input: CreateBookInput{Title: "  Go in Action ", Author: " William Kennedy "}},
		{name: "rejects blank title", repo: &fakeRepository{}, input: CreateBookInput{Title: " ", Author: "Author"}, wantErr: ErrValidation},
		{name: "rejects duplicate ignoring case", repo: &fakeRepository{books: []Book{{Title: "Let's Go", Author: "Alex Edwards"}}}, input: CreateBookInput{Title: " let's go ", Author: "ALEX EDWARDS"}, wantErr: ErrDuplicateBook},
		{name: "returns list failure", repo: &fakeRepository{listErr: repositoryErr}, input: CreateBookInput{Title: "Go", Author: "Author"}, wantErr: repositoryErr},
		{name: "returns create failure", repo: &fakeRepository{createErr: repositoryErr}, input: CreateBookInput{Title: "Go", Author: "Author"}, wantErr: repositoryErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewService(tt.repo).Create(context.Background(), tt.input)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Create() error = %v", err)
				}
				if tt.repo.created != (CreateBookInput{Title: "Go in Action", Author: "William Kennedy"}) {
					t.Fatalf("created input = %+v", tt.repo.created)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceGetMapsMissingBook(t *testing.T) {
	_, err := NewService(&fakeRepository{getErr: ErrBookNotFound}).Get(context.Background(), 99)
	if !errors.Is(err, ErrBookNotFound) {
		t.Fatalf("Get() error = %v, want ErrBookNotFound", err)
	}
}
