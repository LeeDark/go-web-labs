package router

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/http/handlers"
)

func newTestRouter() http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := handlers.NewBooksHandler(books.NewService(books.NewMemoryRepository()), logger)
	return New(handler, logger)
}

func TestHealthReturnsJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	newTestRouter().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health", nil))

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

func TestDeleteThenGetReturnsNotFound(t *testing.T) {
	router := newTestRouter()

	deleteRecorder := httptest.NewRecorder()
	router.ServeHTTP(deleteRecorder, httptest.NewRequest(http.MethodDelete, "/books/1", nil))
	if deleteRecorder.Code != http.StatusNoContent || deleteRecorder.Body.Len() != 0 {
		t.Fatalf("DELETE response = %d %q", deleteRecorder.Code, deleteRecorder.Body.String())
	}

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, httptest.NewRequest(http.MethodGet, "/books/1", nil))
	if getRecorder.Code != http.StatusNotFound {
		t.Fatalf("GET after DELETE status = %d, want %d", getRecorder.Code, http.StatusNotFound)
	}
}

func TestPatchThenGetReturnsUpdatedBook(t *testing.T) {
	router := newTestRouter()

	patchRecorder := httptest.NewRecorder()
	patchRequest := httptest.NewRequest(http.MethodPatch, "/books/1", strings.NewReader(`{"title":"Updated Go"}`))
	patchRequest.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(patchRecorder, patchRequest)
	if patchRecorder.Code != http.StatusOK {
		t.Fatalf("PATCH status = %d, want %d", patchRecorder.Code, http.StatusOK)
	}

	getRecorder := httptest.NewRecorder()
	router.ServeHTTP(getRecorder, httptest.NewRequest(http.MethodGet, "/books/1", nil))
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d", getRecorder.Code, http.StatusOK)
	}
	var response struct {
		Data struct {
			ID     int64  `json:"id"`
			Title  string `json:"title"`
			Author string `json:"author"`
		} `json:"data"`
	}
	if err := json.NewDecoder(getRecorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if response.Data.ID != 1 || response.Data.Title != "Updated Go" || response.Data.Author != "Alan A. A. Donovan" {
		t.Fatalf("GET data = %+v", response.Data)
	}
}
