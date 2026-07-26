package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				app.logError(r, fmt.Errorf("panic recovered: %v", err))
				app.serverErrorResponse(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
