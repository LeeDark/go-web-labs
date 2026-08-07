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
	return &application{logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
}

func TestCreateMovieHandlerRejectsInvalidJSON(t *testing.T) {
	app := newTestApplication()

	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "malformed", body: `{"title":`, want: "badly-formed JSON"},
		{name: "unknown field", body: `{"title":"Go","unknown":true}`, want: "unknown field"},
		{name: "empty", body: "", want: "must not be empty"},
		{name: "multiple values", body: `{} {}`, want: "single JSON value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(tt.body))
			app.createMovieHandler(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q", got)
			}
			var response struct {
				Error string `json:"error"`
			}
			if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if !strings.Contains(response.Error, tt.want) {
				t.Fatalf("error = %q, want substring %q", response.Error, tt.want)
			}
		})
	}
}

func TestCreateMovieHandlerReturnsValidationErrorBeforeDatabaseAccess(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/movies", strings.NewReader(`{"title":"","year":0,"runtime":"0 mins","genres":null}`))

	app.createMovieHandler(recorder, request)

	if recorder.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnprocessableEntity)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
	var response struct {
		Error map[string]string `json:"error"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	for _, field := range []string{"title", "year", "runtime", "genres"} {
		if response.Error[field] == "" {
			t.Fatalf("validation response missing %q: %#v", field, response.Error)
		}
	}
}

func TestRoutesReturnJSONForUnsupportedMethod(t *testing.T) {
	app := newTestApplication()
	recorder := httptest.NewRecorder()

	app.routes().ServeHTTP(recorder, httptest.NewRequest(http.MethodPut, "/v1/movies", nil))

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := recorder.Header().Get("Allow"); got != "OPTIONS, POST" {
		t.Fatalf("Allow = %q, want %q", got, "OPTIONS, POST")
	}
	var response struct {
		Error string `json:"error"`
	}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !strings.Contains(response.Error, "PUT method is not supported") {
		t.Fatalf("error = %q", response.Error)
	}
}
