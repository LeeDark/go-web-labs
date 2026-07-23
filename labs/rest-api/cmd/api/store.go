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

func (s *bookStore) create(book Book) Book {
	s.mu.Lock()
	defer s.mu.Unlock()

	book.ID = s.nextID
	s.nextID++
	s.books[book.ID] = book

	return book
}

func (s *bookStore) update(id int64, book Book) (Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.books[id]
	if !ok {
		return Book{}, false
	}

	book.ID = id
	s.books[id] = book
	return book, true
}

func (s *bookStore) delete(id int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.books[id]
	if !ok {
		return false
	}

	delete(s.books, id)
	return true
}
