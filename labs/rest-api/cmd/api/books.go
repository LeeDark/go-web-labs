package main

import (
	"net/http"
)

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
