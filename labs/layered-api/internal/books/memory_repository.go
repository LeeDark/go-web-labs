package books

import (
	"context"
	"sync"
)

// MemoryRepository is a concurrency-safe learning adapter, not production storage.
type MemoryRepository struct {
	mu     sync.RWMutex
	books  map[int64]Book
	nextID int64
}

func NewMemoryRepository() *MemoryRepository {
	seed := []Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan"},
		{ID: 2, Title: "Let's Go", Author: "Alex Edwards"},
	}

	books := make(map[int64]Book, len(seed))
	for _, book := range seed {
		books[book.ID] = book
	}

	return &MemoryRepository{books: books, nextID: 3}
}

func (r *MemoryRepository) List(_ context.Context) ([]Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Book, 0, len(r.books))
	for id := int64(1); id < r.nextID; id++ {
		if book, ok := r.books[id]; ok {
			result = append(result, book)
		}
	}
	return result, nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id int64) (Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.books[id]
	if !ok {
		return Book{}, ErrBookNotFound
	}
	return book, nil
}

func (r *MemoryRepository) Create(_ context.Context, input CreateBookInput) (Book, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	book := Book{ID: r.nextID, Title: input.Title, Author: input.Author}
	r.books[book.ID] = book
	r.nextID++
	return book, nil
}
