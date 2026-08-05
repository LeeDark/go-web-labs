package main

import (
	"errors"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/app"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/router"
	"github.com/samber/do/v2"
)

func main() {
	addr := flag.String("addr", ":4001", "API server address")
	flag.Parse()

	server, err := newServer(*addr)
	if err != nil {
		log.Printf("build server: %v", err)
		os.Exit(1)
	}

	log.Printf("starting server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}

func newServer(addr string) (*http.Server, error) {
	injector := do.New()
	registerServices(injector, addr)

	return do.Invoke[*http.Server](injector)
}

func registerServices(injector do.Injector, addr string) {
	do.ProvideValue(injector, addr)
	do.Provide(injector, func(do.Injector) (*slog.Logger, error) {
		return app.NewLogger(), nil
	})
	do.Provide(injector, func(do.Injector) (*books.MemoryRepository, error) {
		return books.NewMemoryRepository(), nil
	})
	do.Provide(injector, func(injector do.Injector) (books.BookRepository, error) {
		return do.Invoke[*books.MemoryRepository](injector)
	})
	do.Provide(injector, func(injector do.Injector) (books.Service, error) {
		repository, err := do.Invoke[books.BookRepository](injector)
		if err != nil {
			return nil, err
		}
		return books.NewService(repository), nil
	})
	do.Provide(injector, func(injector do.Injector) (*handlers.BooksHandler, error) {
		service, err := do.Invoke[books.Service](injector)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](injector)
		if err != nil {
			return nil, err
		}
		return handlers.NewBooksHandler(service, logger), nil
	})
	do.Provide(injector, func(injector do.Injector) (http.Handler, error) {
		handler, err := do.Invoke[*handlers.BooksHandler](injector)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](injector)
		if err != nil {
			return nil, err
		}
		return router.New(handler, logger), nil
	})
	do.Provide(injector, func(injector do.Injector) (*http.Server, error) {
		handler, err := do.Invoke[http.Handler](injector)
		if err != nil {
			return nil, err
		}
		logger, err := do.Invoke[*slog.Logger](injector)
		if err != nil {
			return nil, err
		}
		address, err := do.Invoke[string](injector)
		if err != nil {
			return nil, err
		}
		return app.NewServer(address, handler, logger), nil
	})
}
