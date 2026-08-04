package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
)

func TestRoutesHealthReturnsJSON(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := handlers.NewBooksHandler(books.NewService(books.NewMemoryRepository()), logger)
	recorder := httptest.NewRecorder()

	routes(handler, logger).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode health response: %v", err)
	}
	if response["status"] != "available" {
		t.Fatalf("status body = %q", response["status"])
	}
}
