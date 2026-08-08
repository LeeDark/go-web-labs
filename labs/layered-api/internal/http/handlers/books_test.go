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
	updateCalls int
	updateID    int64
	updateInput books.UpdateBookInput
	deleteCalls int
	deleteID    int64
}

func (s *fakeService) List(context.Context) ([]books.Book, error)     { return s.items, s.err }
func (s *fakeService) Get(context.Context, int64) (books.Book, error) { return s.book, s.err }
func (s *fakeService) Create(_ context.Context, input books.CreateBookInput) (books.Book, error) {
	s.createInput = input
	return s.book, s.err
}
func (s *fakeService) Update(_ context.Context, id int64, input books.UpdateBookInput) (books.Book, error) {
	s.updateCalls++
	s.updateID = id
	s.updateInput = input
	return s.book, s.err
}
func (s *fakeService) Delete(_ context.Context, id int64) error {
	s.deleteCalls++
	s.deleteID = id
	return s.err
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
		{name: "internal error", id: "1", service: &fakeService{err: errors.New("database password leaked")}, want: http.StatusInternalServerError, body: `"code":"internal_error"`},
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
			if test.name == "internal error" && strings.Contains(recorder.Body.String(), "password") {
				t.Fatalf("internal error exposed details: %s", recorder.Body.String())
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

func TestBooksHandlerUpdate(t *testing.T) {
	service := &fakeService{book: books.Book{ID: 1, Title: "Updated", Author: "Author"}}
	handler := newTestHandler(service)
	recorder := httptest.NewRecorder()
	request := jsonRequest(`{"title":"Updated"}`)
	request = withRouteID(request, "1")
	handler.Update(recorder, request)

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"title":"Updated"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
	if service.updateInput.Title == nil || *service.updateInput.Title != "Updated" || service.updateInput.Author != nil {
		t.Fatalf("service update input = %+v", service.updateInput)
	}
}

func TestBooksHandlerUpdateRejectsInvalidInput(t *testing.T) {
	handler := newTestHandler(&fakeService{})
	for _, body := range []string{"", `{}`, `null`, `{"title":null}`, `{"unknown":"field"}`, `{"title":"Go"}{}`} {
		recorder := httptest.NewRecorder()
		request := withRouteID(jsonRequest(body), "1")
		handler.Update(recorder, request)
		if recorder.Code != http.StatusBadRequest || !strings.Contains(recorder.Body.String(), `"code":"invalid_request"`) {
			t.Fatalf("body %q response = %d %s", body, recorder.Code, recorder.Body.String())
		}
	}
}

func TestBooksHandlerUpdateRejectsInvalidIDWithoutCallingService(t *testing.T) {
	service := &fakeService{}
	handler := newTestHandler(service)
	recorder := httptest.NewRecorder()
	request := withRouteID(jsonRequest(`{"title":"Updated"}`), "zero")
	handler.Update(recorder, request)

	if recorder.Code != http.StatusBadRequest || !strings.Contains(recorder.Body.String(), `"code":"invalid_id"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
	if service.updateCalls != 0 {
		t.Fatalf("Update called %d times with ID %d", service.updateCalls, service.updateID)
	}
}

func TestBooksHandlerUpdateRejectsUnsupportedMediaType(t *testing.T) {
	handler := newTestHandler(&fakeService{})
	recorder := httptest.NewRecorder()
	request := withRouteID(httptest.NewRequest(http.MethodPatch, "/books/1", strings.NewReader(`{"title":"Updated"}`)), "1")
	request.Header.Set("Content-Type", "text/plain")
	handler.Update(recorder, request)

	if recorder.Code != http.StatusUnsupportedMediaType || !strings.Contains(recorder.Body.String(), `"code":"unsupported_media_type"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestBooksHandlerUpdateMapsApplicationErrors(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
		want int
		code string
	}{
		{name: "validation", err: books.ErrValidation, want: http.StatusUnprocessableEntity, code: "validation_failed"},
		{name: "duplicate", err: books.ErrDuplicateBook, want: http.StatusConflict, code: "duplicate_book"},
		{name: "missing", err: books.ErrBookNotFound, want: http.StatusNotFound, code: "book_not_found"},
	} {
		t.Run(test.name, func(t *testing.T) {
			handler := newTestHandler(&fakeService{err: test.err})
			recorder := httptest.NewRecorder()
			request := withRouteID(jsonRequest(`{"title":"Updated"}`), "1")
			handler.Update(recorder, request)
			if recorder.Code != test.want || !strings.Contains(recorder.Body.String(), `"code":"`+test.code+`"`) {
				t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestBooksHandlerUpdateReturnsSafeInternalError(t *testing.T) {
	handler := newTestHandler(&fakeService{err: errors.New("database password leaked")})
	recorder := httptest.NewRecorder()
	request := withRouteID(jsonRequest(`{"title":"Updated"}`), "1")
	handler.Update(recorder, request)

	if recorder.Code != http.StatusInternalServerError || strings.Contains(recorder.Body.String(), "password") || !strings.Contains(recorder.Body.String(), `"code":"internal_error"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestBooksHandlerDelete(t *testing.T) {
	handler := newTestHandler(&fakeService{})
	recorder := httptest.NewRecorder()
	request := withRouteID(httptest.NewRequest(http.MethodDelete, "/books/1", nil), "1")
	handler.Delete(recorder, request)

	if recorder.Code != http.StatusNoContent || recorder.Body.Len() != 0 {
		t.Fatalf("response = %d %q", recorder.Code, recorder.Body.String())
	}
}

func TestBooksHandlerDeleteMapsMissingBook(t *testing.T) {
	handler := newTestHandler(&fakeService{err: books.ErrBookNotFound})
	recorder := httptest.NewRecorder()
	request := withRouteID(httptest.NewRequest(http.MethodDelete, "/books/99", nil), "99")
	handler.Delete(recorder, request)

	if recorder.Code != http.StatusNotFound || !strings.Contains(recorder.Body.String(), `"code":"book_not_found"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
}

func TestBooksHandlerDeleteRejectsInvalidIDWithoutCallingService(t *testing.T) {
	service := &fakeService{}
	handler := newTestHandler(service)
	recorder := httptest.NewRecorder()
	request := withRouteID(httptest.NewRequest(http.MethodDelete, "/books/zero", nil), "zero")
	handler.Delete(recorder, request)

	if recorder.Code != http.StatusBadRequest || !strings.Contains(recorder.Body.String(), `"code":"invalid_id"`) {
		t.Fatalf("response = %d %s", recorder.Code, recorder.Body.String())
	}
	if service.deleteCalls != 0 {
		t.Fatalf("Delete called %d times with ID %d", service.deleteCalls, service.deleteID)
	}
}

func TestBooksHandlerDeleteReturnsSafeInternalError(t *testing.T) {
	handler := newTestHandler(&fakeService{err: errors.New("database password leaked")})
	recorder := httptest.NewRecorder()
	request := withRouteID(httptest.NewRequest(http.MethodDelete, "/books/1", nil), "1")
	handler.Delete(recorder, request)

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
