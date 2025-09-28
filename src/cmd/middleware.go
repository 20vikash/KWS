package main

import (
	"context"
	"net/http"
)

func (app *Application) IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is authorized.
		isAuthorized := app.SessionManager.GetBool(r.Context(), "isAuthorized")
		if !isAuthorized {
			http.Redirect(w, r, "/kws_signin", http.StatusSeeOther)
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

// Tunnel authorization
func (app *Application) IsTunnelUserAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("secret")
		if secret == "" {
			http.Error(w, "Missing secret header", http.StatusUnauthorized)
			return
		}

		// Check if the tunnel request is authorized
		uid, err := app.Store.InMemory.GetUidFromTunnelSecret(r.Context(), secret)
		ctx := context.WithValue(r.Context(), "uid", uid)
		r = r.WithContext(ctx)
		if err != nil {
			http.Error(w, "Unauthorized request", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
