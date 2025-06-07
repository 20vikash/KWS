package main

import "net/http"

func (app *Application) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is authorized.
		isAuthorized := app.sessionManager.GetBool(r.Context(), "isAuthorized")
		if !isAuthorized {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not authorized to access this page"))
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
