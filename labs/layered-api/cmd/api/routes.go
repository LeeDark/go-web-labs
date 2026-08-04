package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/middleware"
	"github.com/go-chi/chi/v5"
)

func routes(booksHandler *handlers.BooksHandler, logger *slog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer(logger))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "available"})
	})
	router.Get("/books", booksHandler.List)
	router.Get("/books/{id}", booksHandler.Get)
	router.Post("/books", booksHandler.Create)

	return router
}
