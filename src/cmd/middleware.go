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

// Initial 5 burst, then 1 per 10 seconds
func (app *Application) LoginRateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getVisitorLimiter(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many login attempts. Try again later.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
