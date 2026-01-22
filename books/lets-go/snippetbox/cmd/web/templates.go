package main

import "github.com/LeeDark/go-web-labs/books/lets-go/snippetbox/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
