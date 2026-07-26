package main

import (
	"net/http"
	"strings"
)

type apiError struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	env := envelope{"error": apiError{Code: code, Message: message}}

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

func (app *application) invalidJSONResponse(w http.ResponseWriter, r *http.Request) {
	app.badRequestResponse(w, r, "invalid_json", "Request body must contain valid JSON")
}

func (app *application) validationFailedResponse(w http.ResponseWriter, r *http.Request) {
	app.validationFailedWithFieldsResponse(w, r, nil)
}

func (app *application) validationFailedWithFieldsResponse(w http.ResponseWriter, r *http.Request, fields map[string]string) {
	env := envelope{"error": apiError{
		Code:    "validation_failed",
		Message: "Request validation failed",
		Fields:  fields,
	}}

	if err := app.writeJSON(w, http.StatusUnprocessableEntity, env, nil); err != nil {
		app.logError(r, err)
	}
}

func (app *application) requestTooLargeResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusRequestEntityTooLarge, "request_too_large", "Request body must not exceed 1 MiB")
}

func (app *application) unsupportedMediaTypeResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json")
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/health":
		w.Header().Set("Allow", http.MethodGet)
	case r.URL.Path == "/books":
		w.Header().Set("Allow", "GET, POST")
	case strings.HasPrefix(r.URL.Path, "/books/"):
		w.Header().Set("Allow", "GET, PATCH, DELETE")
	}

	app.errorResponse(w, r, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
}

func (app *application) routeNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, "route_not_found", "Route not found")
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusInternalServerError, "internal_error", "An unexpected error occurred")
}
