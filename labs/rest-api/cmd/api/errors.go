package main

import "net/http"

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	env := envelope{"error": apiError{code, message}}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
	}
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, code, message string) {
	app.errorResponse(w, r, http.StatusNotFound, code, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, code, message string) {
	app.errorResponse(w, r, http.StatusBadRequest, code, message)
}
