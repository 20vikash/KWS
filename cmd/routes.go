package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(app *Application) http.Handler {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)

	// Endpoints
	r.Get("/", app.HelloWorld)

	return r
}
