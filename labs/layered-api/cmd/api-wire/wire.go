//go:build wireinject

package main

import (
	"net/http"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/app"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/router"
	"github.com/goforj/wire"
)

func initializeServer(addr string) (*http.Server, error) {
	wire.Build(
		app.NewLogger,
		books.NewMemoryRepository,
		wire.Bind(new(books.BookRepository), new(*books.MemoryRepository)),
		books.NewService,
		handlers.NewBooksHandler,
		router.New,
		app.NewServer,
	)
	return nil, nil
}
