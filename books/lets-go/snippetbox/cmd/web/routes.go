package main

import (
	"net/http"

	"github.com/LeeDark/go-web-labs/books/lets-go/snippetbox/internal/middleware/neutered"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neutered.NewFileSystem(http.Dir(app.cfg.staticDir)))

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return mux
}
