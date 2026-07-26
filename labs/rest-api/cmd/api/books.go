package main

import (
	"errors"
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

func (in *createBookInput) validate() map[string]string {
	fields := make(map[string]string)

	validateBookFields(fields, in.Title, in.Author, in.Description, true, true, true)
	return fields
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

func (in *patchBookInput) validate() map[string]string {
	fields := make(map[string]string)
	if in.Title == nil && in.Author == nil && in.Description == nil {
		fields["body"] = "must include at least one supported field"
		return fields
	}
	if in.Title != nil {
		validateBookFields(fields, *in.Title, "", "", true, false, false)
	}
	if in.Author != nil {
		validateBookFields(fields, "", *in.Author, "", false, true, false)
	}
	if in.Description != nil {
		validateBookFields(fields, "", "", *in.Description, false, false, true)
	}
	return fields
}

func validateBookFields(fields map[string]string, title, author, description string, checkTitle, checkAuthor, checkDescription bool) {
	if checkTitle {
		switch {
		case title == "":
			fields["title"] = "must be provided"
		case len([]rune(title)) > 200:
			fields["title"] = "must not exceed 200 characters"
		}
	}
	if checkAuthor {
		switch {
		case author == "":
			fields["author"] = "must be provided"
		case len([]rune(author)) > 120:
			fields["author"] = "must not exceed 120 characters"
		}
	}
	if checkDescription && len([]rune(description)) > 1000 {
		fields["description"] = "must not exceed 1000 characters"
	}
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

	if err := app.readJSON(w, r, &input); err != nil {
		app.handleJSONError(w, r, err)
		return
	}

	input.normalize()
	if fields := input.validate(); len(fields) > 0 {
		app.validationFailedWithFieldsResponse(w, r, fields)
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
	if err := app.readJSON(w, r, &input); err != nil {
		app.handleJSONError(w, r, err)
		return
	}

	input.normalize()
	if fields := input.validate(); len(fields) > 0 {
		app.validationFailedWithFieldsResponse(w, r, fields)
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

func (app *application) handleJSONError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, errRequestTooLarge):
		app.requestTooLargeResponse(w, r)
	case errors.Is(err, errUnsupportedMediaType):
		app.unsupportedMediaTypeResponse(w, r)
	default:
		app.invalidJSONResponse(w, r)
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
