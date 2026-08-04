package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	items       []books.Book
	book        books.Book
	err         error
	createInput books.CreateBookInput
}

func (s *fakeService) List(context.Context) ([]books.Book, error)     { return s.items, s.err }
func (s *fakeService) Get(context.Context, int64) (books.Book, error) { return s.book, s.err }
func (s *fakeService) Create(_ context.Context, input books.CreateBookInput) (books.Book, error) {
	s.createInput = input
	return s.book, s.err
}

func newTestHandler(service books.Service) *BooksHandler {
	return NewBooksHandler(service, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestBooksHandlerList(t *testing.T) {
	handler := newTestHandler(&fakeService{items: []books.Book{{ID: 1, Title: "Go", Author: "Author"}}})
	recorder := httptest.NewRecorder()
	handler.List(recorder, httptest.NewRequest(http.MethodGet, "/books", nil))

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"data"`) || !strings.Contains(recorder.Body.String(), `"title":"Go"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
}

func TestBooksHandlerGet(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		service *fakeService
		want    int
		body    string
	}{
		{name: "existing", id: "1", service: &fakeService{book: books.Book{ID: 1, Title: "Go", Author: "Author"}}, want: http.StatusOK, body: `"title":"Go"`},
		{name: "missing", id: "99", service: &fakeService{err: books.ErrBookNotFound}, want: http.StatusNotFound, body: `"code":"book_not_found"`},
		{name: "invalid ID", id: "zero", service: &fakeService{}, want: http.StatusBadRequest, body: `"code":"invalid_id"`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := newTestHandler(test.service)
			recorder := httptest.NewRecorder()
			request := withRouteID(httptest.NewRequest(http.MethodGet, "/books/"+test.id, nil), test.id)
			handler.Get(recorder, request)
			if recorder.Code != test.want || !strings.Contains(recorder.Body.String(), test.body) {
				t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestBooksHandlerCreate(t *testing.T) {
	service := &fakeService{book: books.Book{ID: 3, Title: "Go", Author: "Author"}}
	handler := newTestHandler(service)
	recorder := httptest.NewRecorder()
	handler.Create(recorder, jsonRequest(`{"title":"Go","author":"Author"}`))

	if recorder.Code != http.StatusCreated || recorder.Header().Get("Location") != "/books/3" {
		t.Fatalf("response = %d location=%q", recorder.Code, recorder.Header().Get("Location"))
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
	if service.createInput != (books.CreateBookInput{Title: "Go", Author: "Author"}) {
		t.Fatalf("service input = %+v", service.createInput)
	}
	if !strings.Contains(recorder.Body.String(), `"data"`) {
		t.Fatalf("response does not contain data envelope: %s", recorder.Body.String())
	}
}

func TestBooksHandlerCreateMapsValidationAndDuplicate(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
		want int
		code string
	}{
		{name: "validation", err: books.ErrValidation, want: http.StatusUnprocessableEntity, code: "validation_failed"},
		{name: "duplicate", err: books.ErrDuplicateBook, want: http.StatusConflict, code: "duplicate_book"},
	} {
		t.Run(test.name, func(t *testing.T) {
			handler := newTestHandler(&fakeService{err: test.err})
			recorder := httptest.NewRecorder()
			handler.Create(recorder, jsonRequest(`{"title":"Go","author":"Author"}`))
			if recorder.Code != test.want || !strings.Contains(recorder.Body.String(), `"code":"`+test.code+`"`) {
				t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestBooksHandlerRejectsInvalidInput(t *testing.T) {
	handler := newTestHandler(&fakeService{})
	for _, body := range []string{`{`, `{"title":"Go","unknown":true}`, `{"title":"Go"}{}`} {
		recorder := httptest.NewRecorder()
		handler.Create(recorder, jsonRequest(body))
		if recorder.Code != http.StatusBadRequest || !strings.Contains(recorder.Body.String(), `"code":"invalid_request"`) {
			t.Fatalf("body %q status = %d", body, recorder.Code)
		}
	}
}

func TestBooksHandlerCreateRejectsUnsupportedMediaType(t *testing.T) {
	handler := newTestHandler(&fakeService{})
	for _, contentType := range []string{"", "text/plain", "application/json; charset="} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(`{"title":"Go","author":"Author"}`))
		if contentType != "" {
			request.Header.Set("Content-Type", contentType)
		}
		handler.Create(recorder, request)
		if recorder.Code != http.StatusUnsupportedMediaType || !strings.Contains(recorder.Body.String(), `"code":"unsupported_media_type"`) {
			t.Fatalf("Content-Type %q response = %d %s", contentType, recorder.Code, recorder.Body.String())
		}
	}
}

func TestBooksHandlerReturnsSafeInternalError(t *testing.T) {
	handler := newTestHandler(&fakeService{err: errors.New("database password leaked")})
	recorder := httptest.NewRecorder()
	handler.List(recorder, httptest.NewRequest(http.MethodGet, "/books", nil))
	if recorder.Code != http.StatusInternalServerError || strings.Contains(recorder.Body.String(), "password") || !strings.Contains(recorder.Body.String(), `"code":"internal_error"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
}

func jsonRequest(body string) *http.Request {
	request := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	return request
}

func withRouteID(request *http.Request, id string) *http.Request {
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", id)
	return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
}
