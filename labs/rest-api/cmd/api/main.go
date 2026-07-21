package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type config struct {
	addr string
}

type application struct {
	config    config
	logger    *slog.Logger
	bookStore *bookStore
}

func main() {
	var cfg config

	flag.StringVar(&cfg.addr, "addr", ":4000", "API server address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := application{
		config:    cfg,
		logger:    logger,
		bookStore: newBookStore(),
	}

	srv := &http.Server{
		Addr:              cfg.addr,
		Handler:           app.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       time.Minute,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
