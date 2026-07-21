package main

import "net/http"

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"data": map[string]string{
			"status": "available",
		},
	}

	if err := app.writeJSON(w, http.StatusOK, data, nil); err != nil {
		app.logger.Error("writing health response",
			"error", err,
			"method", r.Method,
			"uri", r.URL.RequestURI(),
		)
	}
}
