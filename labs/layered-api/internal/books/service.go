package books

import (
	"context"
	"errors"
	"strings"
	"sync"
)

var (
	ErrValidation    = errors.New("book validation failed")
	ErrDuplicateBook = errors.New("duplicate book")
	ErrBookNotFound  = errors.New("book not found")
)

// Service exposes book use cases without HTTP concerns.
type Service interface {
	List(ctx context.Context) ([]Book, error)
	Get(ctx context.Context, id int64) (Book, error)
	Create(ctx context.Context, input CreateBookInput) (Book, error)
}

type bookService struct {
	repository BookRepository
	createMu   sync.Mutex
}

func NewService(repository BookRepository) Service {
	return &bookService{repository: repository}
}

func (s *bookService) List(ctx context.Context) ([]Book, error) {
	return s.repository.List(ctx)
}

func (s *bookService) Get(ctx context.Context, id int64) (Book, error) {
	book, err := s.repository.GetByID(ctx, id)
	if errors.Is(err, ErrBookNotFound) {
		return Book{}, ErrBookNotFound
	}
	return book, err
}

func (s *bookService) Create(ctx context.Context, input CreateBookInput) (Book, error) {
	// The duplicate check and create must be one critical section for this in-memory lab.
	s.createMu.Lock()
	defer s.createMu.Unlock()

	input.Title = strings.TrimSpace(input.Title)
	input.Author = strings.TrimSpace(input.Author)
	if input.Title == "" || input.Author == "" {
		return Book{}, ErrValidation
	}

	books, err := s.repository.List(ctx)
	if err != nil {
		return Book{}, err
	}
	for _, book := range books {
		if strings.EqualFold(book.Title, input.Title) && strings.EqualFold(book.Author, input.Author) {
			return Book{}, ErrDuplicateBook
		}
	}

	return s.repository.Create(ctx, input)
}
