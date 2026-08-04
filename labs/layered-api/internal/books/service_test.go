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
	updateErr error
	deleteErr error
	created   CreateBookInput
	updated   Book
	deletedID int64
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
func (r *fakeRepository) Update(_ context.Context, book Book) (Book, error) {
	r.updated = book
	if r.updateErr != nil {
		return Book{}, r.updateErr
	}
	return book, nil
}
func (r *fakeRepository) Delete(_ context.Context, id int64) error {
	r.deletedID = id
	return r.deleteErr
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

func TestServiceUpdate(t *testing.T) {
	title := "  Updated title "
	blank := " "
	duplicateTitle := "other book"
	duplicateAuthor := "OTHER AUTHOR"
	repositoryErr := errors.New("database unavailable")

	tests := []struct {
		name    string
		repo    *fakeRepository
		input   UpdateBookInput
		want    Book
		wantErr error
	}{
		{
			name:  "updates one field and preserves the other",
			repo:  &fakeRepository{getBook: Book{ID: 1, Title: "Original title", Author: "Original author"}},
			input: UpdateBookInput{Title: &title},
			want:  Book{ID: 1, Title: "Updated title", Author: "Original author"},
		},
		{
			name:    "rejects blank value after merge",
			repo:    &fakeRepository{getBook: Book{ID: 1, Title: "Original title", Author: "Original author"}},
			input:   UpdateBookInput{Title: &blank},
			wantErr: ErrValidation,
		},
		{
			name: "rejects duplicate after merge",
			repo: &fakeRepository{
				getBook: Book{ID: 1, Title: "Original title", Author: "Original author"},
				books:   []Book{{ID: 1, Title: "Original title", Author: "Original author"}, {ID: 2, Title: "Other Book", Author: "Other Author"}},
			},
			input:   UpdateBookInput{Title: &duplicateTitle, Author: &duplicateAuthor},
			wantErr: ErrDuplicateBook,
		},
		{
			name:    "maps missing book",
			repo:    &fakeRepository{getErr: ErrBookNotFound},
			input:   UpdateBookInput{Title: &title},
			wantErr: ErrBookNotFound,
		},
		{
			name:    "returns repository update failure",
			repo:    &fakeRepository{getBook: Book{ID: 1, Title: "Original title", Author: "Original author"}, updateErr: repositoryErr},
			input:   UpdateBookInput{Title: &title},
			wantErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewService(tt.repo).Update(context.Background(), 1, tt.input)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Update() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Update() error = %v", err)
			}
			if got != tt.want || tt.repo.updated != tt.want {
				t.Fatalf("Update() book = %+v, repository update = %+v, want %+v", got, tt.repo.updated, tt.want)
			}
		})
	}
}

func TestServiceDelete(t *testing.T) {
	repositoryErr := errors.New("database unavailable")
	tests := []struct {
		name    string
		repo    *fakeRepository
		wantErr error
	}{
		{name: "deletes existing book", repo: &fakeRepository{}},
		{name: "maps missing book", repo: &fakeRepository{deleteErr: ErrBookNotFound}, wantErr: ErrBookNotFound},
		{name: "returns repository failure", repo: &fakeRepository{deleteErr: repositoryErr}, wantErr: repositoryErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewService(tt.repo).Delete(context.Background(), 7)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Delete() error = %v, want %v", err, tt.wantErr)
			}
			if tt.repo.deletedID != 7 {
				t.Fatalf("repository deleted ID = %d, want 7", tt.repo.deletedID)
			}
		})
	}
}
