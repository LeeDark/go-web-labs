package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/LeeDark/go-web-labs/books/lets-go/snippetbox/internal/middleware/neutered"
)

type config struct {
	addr      string
	staticDir string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(neutered.NewFileSystem(http.Dir(cfg.staticDir)))

	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Printf("starting server on %s", cfg.addr)

	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
}
