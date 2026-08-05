package main

import (
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4001", "API server address")
	flag.Parse()

	server, err := initializeServer(*addr)
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
