package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)

	// Handlers
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// HTTP server
	http.ListenAndServe(":8080", r)
}
