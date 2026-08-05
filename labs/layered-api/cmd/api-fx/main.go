package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/app"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/router"
	"go.uber.org/fx"
)

func main() {
	addr := flag.String("addr", ":4001", "API server address")
	flag.Parse()

	newApplication(*addr).Run()
}

func newApplication(addr string) *fx.App {
	return fx.New(applicationOptions(addr)...)
}

func applicationOptions(addr string) []fx.Option {
	return []fx.Option{
		fx.NopLogger,
		fx.Supply(addr),
		fx.Provide(
			app.NewLogger,
			fx.Annotate(books.NewMemoryRepository, fx.As(new(books.BookRepository))),
			books.NewService,
			handlers.NewBooksHandler,
			router.New,
			app.NewServer,
		),
		fx.Invoke(registerServer),
	}
}

func registerServer(lifecycle fx.Lifecycle, server *http.Server, logger *slog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("starting server", "addr", server.Addr)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error("server stopped", "error", err)
				}
			}()
			return nil
		},
		OnStop: server.Shutdown,
	})
}
