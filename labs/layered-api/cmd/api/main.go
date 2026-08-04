package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
)

func main() {
	addr := flag.String("addr", ":4001", "API server address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repository := books.NewMemoryRepository()
	service := books.NewService(repository)
	booksHandler := handlers.NewBooksHandler(service, logger)

	server := &http.Server{
		Addr:              *addr,
		Handler:           routes(booksHandler, logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
