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
	router.Post("/books", app.createBookHandler)
	router.Patch("/books/{id}", app.patchBookHandler)
	router.Delete("/books/{id}", app.deleteBookHandler)

	return router
}
