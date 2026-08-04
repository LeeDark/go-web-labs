package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"strconv"

	"github.com/LeeDark/go-web-labs/labs/layered-api/internal/books"
	"github.com/go-chi/chi/v5"
)

type BooksHandler struct {
	service books.Service
	logger  *slog.Logger
}

func NewBooksHandler(service books.Service, logger *slog.Logger) *BooksHandler {
	return &BooksHandler{service: service, logger: logger}
}

type createBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type bookResponse struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (h *BooksHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.List(r.Context())
	if err != nil {
		h.internalError(w, err)
		return
	}

	response := make([]bookResponse, 0, len(items))
	for _, item := range items {
		response = append(response, mapBook(item))
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": response})
}

func (h *BooksHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, chi.URLParam(r, "id"))
	if !ok {
		return
	}

	book, err := h.service.Get(r.Context(), id)
	if errors.Is(err, books.ErrBookNotFound) {
		writeError(w, http.StatusNotFound, "book_not_found", "Book not found")
		return
	}
	if err != nil {
		h.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": mapBook(book)})
}

func (h *BooksHandler) Create(w http.ResponseWriter, r *http.Request) {
	if !hasJSONContentType(r) {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json")
		return
	}

	var request createBookRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Request body must be a valid JSON object")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "invalid_request", "Request body must contain one JSON object")
		return
	}

	book, err := h.service.Create(r.Context(), books.CreateBookInput{Title: request.Title, Author: request.Author})
	switch {
	case errors.Is(err, books.ErrValidation):
		writeError(w, http.StatusUnprocessableEntity, "validation_failed", "Title and author are required")
		return
	case errors.Is(err, books.ErrDuplicateBook):
		writeError(w, http.StatusConflict, "duplicate_book", "A book with this title and author already exists")
		return
	case err != nil:
		h.internalError(w, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/books/%d", book.ID))
	writeJSON(w, http.StatusCreated, map[string]any{"data": mapBook(book)})
}

func hasJSONContentType(r *http.Request) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	return err == nil && mediaType == "application/json"
}

func (h *BooksHandler) internalError(w http.ResponseWriter, err error) {
	h.logger.Error("book request failed", "error", err)
	writeError(w, http.StatusInternalServerError, "internal_error", "Internal server error")
}

func parseID(w http.ResponseWriter, raw string) (int64, bool) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid_id", "Book ID must be a positive integer")
		return 0, false
	}
	return id, true
}

func mapBook(book books.Book) bookResponse {
	return bookResponse{ID: book.ID, Title: book.Title, Author: book.Author}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	response := errorResponse{}
	response.Error.Code = code
	response.Error.Message = message
	writeJSON(w, status, response)
}
