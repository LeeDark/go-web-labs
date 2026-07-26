package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestApplication() *application {
	return &application{
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		bookStore: newBookStore(),
	}
}

func newJSONRequest(t *testing.T, method, target, body string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeResponse(t *testing.T, recorder *httptest.ResponseRecorder) envelope {
	t.Helper()

	var response envelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}

func TestListBooksHandler(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/books", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}

	response := decodeResponse(t, recorder)
	books, ok := response["data"].([]any)
	if !ok || len(books) != 1 {
		t.Fatalf("data = %#v, want one book", response["data"])
	}
}

func TestShowBookHandler(t *testing.T) {
	app := newTestApplication()

	t.Run("existing", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/books/1", nil))

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
		response := decodeResponse(t, recorder)
		book := response["data"].(map[string]any)
		if book["id"] != float64(1) {
			t.Fatalf("id = %#v, want 1", book["id"])
		}
	})

	t.Run("missing", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/books/999", nil))

		if recorder.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
		}
		response := decodeResponse(t, recorder)
		err := response["error"].(map[string]any)
		if err["code"] != "book_not_found" {
			t.Fatalf("error.code = %#v, want book_not_found", err["code"])
		}
	})
}

func TestCreateBookHandler(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, newJSONRequest(t, http.MethodPost, "/books", `{"title":" Dune ","author":" Frank Herbert ","description":" A science fiction novel. "}`))

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := recorder.Header().Get("Location"); got != "/books/2" {
		t.Fatalf("Location = %q, want /books/2", got)
	}
	response := decodeResponse(t, recorder)
	book := response["data"].(map[string]any)
	if book["title"] != "Dune" || book["author"] != "Frank Herbert" {
		t.Fatalf("book = %#v, want normalized values", book)
	}
}

func TestCreateBookValidation(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, newJSONRequest(t, http.MethodPost, "/books", `{"title":" ","author":""}`))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	response := decodeResponse(t, recorder)
	err := response["error"].(map[string]any)
	if err["code"] != "validation_failed" {
		t.Fatalf("error.code = %#v, want validation_failed", err["code"])
	}
	fields := err["fields"].(map[string]any)
	if fields["title"] != "must be provided" || fields["author"] != "must be provided" {
		t.Fatalf("fields = %#v, want required title and author errors", fields)
	}
}

func TestPatchBookHandler(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, newJSONRequest(t, http.MethodPatch, "/books/1", `{"title":"A new title"}`))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	response := decodeResponse(t, recorder)
	book := response["data"].(map[string]any)
	if book["title"] != "A new title" || book["author"] != "Ursula K. Le Guin" || book["description"] != "A science fiction novel." {
		t.Fatalf("book = %#v, want only title updated", book)
	}
}

func TestDeleteBookHandler(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodDelete, "/books/1", nil))

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNoContent)
	}
	if recorder.Body.Len() != 0 {
		t.Fatalf("body = %q, want empty", recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/books/1", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("read after delete status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}

func TestHTTPBoundaryErrors(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name    string
		request *http.Request
		status  int
		code    string
	}{
		{
			name:    "missing content type",
			request: httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Dune","author":"Frank Herbert"}`)),
			status:  http.StatusUnsupportedMediaType,
			code:    "unsupported_media_type",
		},
		{
			name:    "unknown field",
			request: newJSONRequest(t, http.MethodPost, "/books", `{"title":"Dune","author":"Frank Herbert","id":2}`),
			status:  http.StatusBadRequest,
			code:    "invalid_json",
		},
		{
			name:    "null field",
			request: newJSONRequest(t, http.MethodPatch, "/books/1", `{"title":null}`),
			status:  http.StatusBadRequest,
			code:    "invalid_json",
		},
		{
			name:    "too large",
			request: newJSONRequest(t, http.MethodPost, "/books", `{"title":"`+strings.Repeat("x", maxRequestBodyBytes)+`","author":"Frank Herbert"}`),
			status:  http.StatusRequestEntityTooLarge,
			code:    "request_too_large",
		},
		{
			name:    "method not allowed",
			request: httptest.NewRequest(http.MethodPut, "/books", nil),
			status:  http.StatusMethodNotAllowed,
			code:    "method_not_allowed",
		},
		{
			name:    "route not found",
			request: httptest.NewRequest(http.MethodGet, "/missing", nil),
			status:  http.StatusNotFound,
			code:    "route_not_found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			app.routes().ServeHTTP(recorder, test.request)

			if recorder.Code != test.status {
				t.Fatalf("status = %d, want %d", recorder.Code, test.status)
			}
			response := decodeResponse(t, recorder)
			err := response["error"].(map[string]any)
			if err["code"] != test.code {
				t.Fatalf("error.code = %#v, want %s", err["code"], test.code)
			}
			if test.status == http.StatusMethodNotAllowed && recorder.Header().Get("Allow") != "GET, POST" {
				t.Fatalf("Allow = %q, want GET, POST", recorder.Header().Get("Allow"))
			}
		})
	}

	t.Run("recovery returns safe internal error", func(t *testing.T) {
		handler := app.recoverPanic(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic("unexpected failure")
		}))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/books", nil))

		if recorder.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusInternalServerError)
		}
		response := decodeResponse(t, recorder)
		err := response["error"].(map[string]any)
		if err["code"] != "internal_error" {
			t.Fatalf("error.code = %#v, want internal_error", err["code"])
		}
	})
}
