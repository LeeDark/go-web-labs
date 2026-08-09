package main

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"
)

type fakeServer struct {
	listen   func() error
	shutdown func(context.Context) error
}

func (s fakeServer) ListenAndServe() error {
	return s.listen()
}

func (s fakeServer) Shutdown(ctx context.Context) error {
	return s.shutdown(ctx)
}

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRunServerContextCancellationCallsShutdownWithDeadline(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdownCalled := make(chan context.Context, 1)
	listenResult := make(chan error, 1)
	server := fakeServer{
		listen: func() error {
			return <-listenResult
		},
		shutdown: func(ctx context.Context) error {
			shutdownCalled <- ctx
			listenResult <- http.ErrServerClosed
			return nil
		},
	}

	cancel()
	err := runServer(ctx, server, discardLogger(), time.Second)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	shutdownCtx := <-shutdownCalled
	if _, ok := shutdownCtx.Deadline(); !ok {
		t.Fatal("want shutdown context to have a deadline")
	}
}

func TestRunServerListenAndServeClosed(t *testing.T) {
	server := fakeServer{
		listen: func() error { return http.ErrServerClosed },
		shutdown: func(context.Context) error {
			t.Fatal("shutdown should not be called")
			return nil
		},
	}

	if err := runServer(context.Background(), server, discardLogger(), time.Second); err != nil {
		t.Fatalf("want no error, got %v", err)
	}
}

func TestRunServerUnexpectedListenerError(t *testing.T) {
	wantErr := errors.New("listener failed")
	server := fakeServer{
		listen: func() error { return wantErr },
		shutdown: func(context.Context) error {
			t.Fatal("shutdown should not be called")
			return nil
		},
	}

	if err := runServer(context.Background(), server, discardLogger(), time.Second); !errors.Is(err, wantErr) {
		t.Fatalf("want %v, got %v", wantErr, err)
	}
}

func TestRunServerUnexpectedListenerErrorAfterSuccessfulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wantErr := errors.New("listener failed after shutdown")
	listenResult := make(chan error, 1)
	server := fakeServer{
		listen: func() error { return <-listenResult },
		shutdown: func(context.Context) error {
			listenResult <- wantErr
			return nil
		},
	}

	cancel()
	if err := runServer(ctx, server, discardLogger(), time.Second); !errors.Is(err, wantErr) {
		t.Fatalf("want %v, got %v", wantErr, err)
	}
}

func TestRunServerShutdownError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wantErr := errors.New("shutdown failed")
	listenResult := make(chan error, 1)
	server := fakeServer{
		listen: func() error { return <-listenResult },
		shutdown: func(context.Context) error {
			listenResult <- http.ErrServerClosed
			return wantErr
		},
	}

	cancel()
	if err := runServer(ctx, server, discardLogger(), time.Second); !errors.Is(err, wantErr) {
		t.Fatalf("want %v, got %v", wantErr, err)
	}
}
