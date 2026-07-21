package main

import "sync"

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
