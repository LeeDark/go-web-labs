package books

import "context"

// BookRepository is the persistence boundary required by Service.
type BookRepository interface {
	List(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id int64) (Book, error)
	Create(ctx context.Context, input CreateBookInput) (Book, error)
}
