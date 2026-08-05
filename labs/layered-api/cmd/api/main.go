package main

import (
	"errors"
	"flag"
	"net/http"
	"os"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/app"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/router"
)

func main() {
	addr := flag.String("addr", ":4001", "API server address")
	flag.Parse()

	logger := app.NewLogger()
	repository := books.NewMemoryRepository()
	service := books.NewService(repository)
	booksHandler := handlers.NewBooksHandler(service, logger)
	server := app.NewServer(*addr, router.New(booksHandler, logger), logger)

	logger.Info("starting server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
