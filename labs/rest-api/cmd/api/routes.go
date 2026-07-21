package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Get("/health", app.healthHandler)
	router.Get("/books", app.listBooksHandler)
	router.Get("/books/{id}", app.showBookHandler)

	return router
}
