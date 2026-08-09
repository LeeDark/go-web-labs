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

func assertBookData(t *testing.T, data any, want map[string]any) {
	t.Helper()

	book, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("data = %#v, want book object", data)
	}
	if len(book) != len(want) {
		t.Fatalf("book fields = %#v, want %#v", book, want)
	}
	for field, wantValue := range want {
		if got := book[field]; got != wantValue {
			t.Fatalf("book[%q] = %#v, want %#v", field, got, wantValue)
		}
	}
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
	assertBookData(t, books[0], map[string]any{
		"id":          float64(1),
		"title":       "The Left Hand of Darkness",
		"author":      "Ursula K. Le Guin",
		"description": "A science fiction novel.",
	})
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
		assertBookData(t, response["data"], map[string]any{
			"id":          float64(1),
			"title":       "The Left Hand of Darkness",
			"author":      "Ursula K. Le Guin",
			"description": "A science fiction novel.",
		})
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

	app.routes().ServeHTTP(recorder,
		newJSONRequest(t, http.MethodPost,
			"/books",
			`{"title":" Dune ","author":" Frank Herbert ","description":" A science fiction novel. "}`,
		),
	)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := recorder.Header().Get("Location"); got != "/books/2" {
		t.Fatalf("Location = %q, want /books/2", got)
	}
	response := decodeResponse(t, recorder)
	assertBookData(t, response["data"], map[string]any{
		"id":          float64(2),
		"title":       "Dune",
		"author":      "Frank Herbert",
		"description": "A science fiction novel.",
	})
}

func TestCreateBookValidation(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		fields map[string]string
	}{
		{
			name: "blank required fields",
			body: `{"title":" ","author":""}`,
			fields: map[string]string{
				"title":  "must be provided",
				"author": "must be provided",
			},
		},
		{
			name: "title too long",
			body: `{"title":"` + strings.Repeat("界", 201) + `","author":"Frank Herbert"}`,
			fields: map[string]string{
				"title": "must not exceed 200 characters",
			},
		},
		{
			name: "author too long",
			body: `{"title":"Dune","author":"` + strings.Repeat("界", 121) + `"}`,
			fields: map[string]string{
				"author": "must not exceed 120 characters",
			},
		},
		{
			name: "description too long",
			body: `{"title":"Dune","author":"Frank Herbert","description":"` + strings.Repeat("界", 1001) + `"}`,
			fields: map[string]string{
				"description": "must not exceed 1000 characters",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app := newTestApplication()
			recorder := httptest.NewRecorder()

			app.routes().ServeHTTP(recorder, newJSONRequest(t, http.MethodPost, "/books", test.body))

			if recorder.Code != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
			}
			response := decodeResponse(t, recorder)
			apiErr := response["error"].(map[string]any)
			if apiErr["code"] != "validation_failed" {
				t.Fatalf("error.code = %#v, want validation_failed", apiErr["code"])
			}
			fields := apiErr["fields"].(map[string]any)
			for field, message := range test.fields {
				if fields[field] != message {
					t.Fatalf("fields[%q] = %#v, want %q", field, fields[field], message)
				}
			}
		})
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
	assertBookData(t, response["data"], map[string]any{
		"id":          float64(1),
		"title":       "A new title",
		"author":      "Ursula K. Le Guin",
		"description": "A science fiction novel.",
	})

	recorder = httptest.NewRecorder()
	app.routes().ServeHTTP(recorder, newJSONRequest(t, http.MethodPatch, "/books/1", `{"title":" "}`))

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("invalid patch status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("invalid patch Content-Type = %q, want application/json", got)
	}
	response = decodeResponse(t, recorder)
	apiErr := response["error"].(map[string]any)
	if apiErr["code"] != "validation_failed" {
		t.Fatalf("error.code = %#v, want validation_failed", apiErr["code"])
	}

	recorder = httptest.NewRecorder()
	app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/books/1", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("read after invalid patch status = %d, want %d", recorder.Code, http.StatusOK)
	}
	response = decodeResponse(t, recorder)
	assertBookData(t, response["data"], map[string]any{
		"id":          float64(1),
		"title":       "A new title",
		"author":      "Ursula K. Le Guin",
		"description": "A science fiction novel.",
	})
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
		allow   string
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
			name:    "top-level null",
			request: newJSONRequest(t, http.MethodPost, "/books", `null`),
			status:  http.StatusBadRequest,
			code:    "invalid_json",
		},
		{
			name:    "trailing JSON value",
			request: newJSONRequest(t, http.MethodPost, "/books", `{"title":"Dune","author":"Frank Herbert"} {}`),
			status:  http.StatusBadRequest,
			code:    "invalid_json",
		},
		{
			name:    "wrong top-level type",
			request: newJSONRequest(t, http.MethodPost, "/books", `[]`),
			status:  http.StatusBadRequest,
			code:    "invalid_json",
		},
		{
			name:    "empty body",
			request: newJSONRequest(t, http.MethodPost, "/books", ``),
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
			allow:   "GET, POST",
		},
		{
			name:    "method not allowed for health",
			request: httptest.NewRequest(http.MethodPost, "/health", nil),
			status:  http.StatusMethodNotAllowed,
			code:    "method_not_allowed",
			allow:   "GET",
		},
		{
			name:    "method not allowed for book item",
			request: httptest.NewRequest(http.MethodPost, "/books/1", nil),
			status:  http.StatusMethodNotAllowed,
			code:    "method_not_allowed",
			allow:   "GET, PATCH, DELETE",
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
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}
			response := decodeResponse(t, recorder)
			err := response["error"].(map[string]any)
			if err["code"] != test.code {
				t.Fatalf("error.code = %#v, want %s", err["code"], test.code)
			}
			if test.allow != "" && recorder.Header().Get("Allow") != test.allow {
				t.Fatalf("Allow = %q, want %q", recorder.Header().Get("Allow"), test.allow)
			}
			if test.name == "too large" {
				books := app.bookStore.list()
				if len(books) != 1 || books[0].ID != 1 {
					t.Fatalf("store after oversized request = %#v, want unchanged seed", books)
				}
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
