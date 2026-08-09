package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.env)

	return runServer(ctx, srv, app.logger, 30*time.Second)
}

type serverLifecycle interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

func runServer(ctx context.Context, srv serverLifecycle, logger *slog.Logger, shutdownTimeout time.Duration) error {
	serverError := make(chan error, 1)
	go func() {
		serverError <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverError:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err

	case <-ctx.Done():
		logger.Info("shutdown signal received", "reason", ctx.Err())

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}

		if err := <-serverError; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		logger.Info("stopped server")
		return nil
	}
}
