package main

import (
	"sort"
	"sync"
)

type Book struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

type bookStore struct {
	mu     sync.RWMutex
	books  map[int64]Book
	nextID int64
}

func newBookStore() *bookStore {
	seed := Book{
		ID:          1,
		Title:       "The Left Hand of Darkness",
		Author:      "Ursula K. Le Guin",
		Description: "A science fiction novel.",
	}

	return &bookStore{
		books:  map[int64]Book{seed.ID: seed},
		nextID: 2,
	}
}

func (s *bookStore) get(id int64) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, ok := s.books[id]
	return book, ok
}

func (s *bookStore) list() []Book {
	s.mu.RLock()

	books := make([]Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	s.mu.RUnlock()

	sort.Slice(books, func(i, j int) bool {
		return books[i].ID < books[j].ID
	})

	return books
}
