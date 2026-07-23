package main

import (
	"fmt"
	"net/http"
	"strings"
)

type createBookInput struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
}

func (in *createBookInput) normalize() {
	in.Title = strings.TrimSpace(in.Title)
	in.Author = strings.TrimSpace(in.Author)
	in.Description = strings.TrimSpace(in.Description)
}

func (in *createBookInput) valid() bool {
	return in.Title != "" && in.Author != ""
}

type patchBookInput struct {
	Title       *string `json:"title"`
	Author      *string `json:"author"`
	Description *string `json:"description"`
}

func (in *patchBookInput) normalize() {
	if in.Title != nil {
		*in.Title = strings.TrimSpace(*in.Title)
	}
	if in.Author != nil {
		*in.Author = strings.TrimSpace(*in.Author)
	}
	if in.Description != nil {
		*in.Description = strings.TrimSpace(*in.Description)
	}
}

func (in *patchBookInput) valid() bool {
	if in.Title == nil && in.Author == nil && in.Description == nil {
		return false
	}
	if in.Title != nil && *in.Title == "" {
		return false
	}
	if in.Author != nil && *in.Author == "" {
		return false
	}
	return true
}

func (in *patchBookInput) applyTo(book *Book) {
	if in.Title != nil {
		book.Title = *in.Title
	}
	if in.Author != nil {
		book.Author = *in.Author
	}
	if in.Description != nil {
		book.Description = *in.Description
	}
}

func (app *application) listBooksHandler(w http.ResponseWriter, r *http.Request) {
	books := app.bookStore.list()

	err := app.writeJSON(w, http.StatusOK, envelope{"data": books}, nil)
	if err != nil {
		app.logError(r, err)
	}
}

func (app *application) showBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "invalid_id", "Book ID must be a positive integer")
		return
	}

	book, ok := app.bookStore.get(id)
	if !ok {
		app.notFoundResponse(w, r, "book_not_found", "Book not found")
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": book}, nil)
	if err != nil {
		app.logError(r, err)
	}
}

func (app *application) createBookHandler(w http.ResponseWriter, r *http.Request) {
	var input createBookInput

	if err := app.readJSON(r, &input); err != nil {
		app.invalidJSONResponse(w, r)
		return
	}

	input.normalize()
	if !input.valid() {
		app.validationFailedResponse(w, r)
		return
	}

	book := app.bookStore.create(Book{
		Title:       input.Title,
		Author:      input.Author,
		Description: input.Description,
	})

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/books/%d", book.ID))

	if err := app.writeJSON(w, http.StatusCreated, envelope{"data": book}, headers); err != nil {
		app.logError(r, err)
	}
}

func (app *application) patchBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "invalid_id", "Book ID must be a positive integer")
		return
	}

	book, ok := app.bookStore.get(id)
	if !ok {
		app.notFoundResponse(w, r, "book_not_found", "Book not found")
		return
	}

	var input patchBookInput
	if err := app.readJSON(r, &input); err != nil {
		app.invalidJSONResponse(w, r)
		return
	}

	input.normalize()
	if !input.valid() {
		app.validationFailedResponse(w, r)
		return
	}

	input.applyTo(&book)

	book, ok = app.bookStore.update(id, book)
	if !ok {
		app.notFoundResponse(w, r, "book_not_found", "Book not found")
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"data": book}, nil); err != nil {
		app.logError(r, err)
	}
}

func (app *application) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, "invalid_id", "Book ID must be a positive integer")
		return
	}

	if deleted := app.bookStore.delete(id); !deleted {
		app.notFoundResponse(w, r, "book_not_found", "Book not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
