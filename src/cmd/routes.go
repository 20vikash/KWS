package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(app *Application) http.Handler {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(sessionManager.LoadAndSave)
	r.Use(app.LoginRateLimitMiddleware)

	// Define a sub-router for protected routes
	r.Group(func(protected chi.Router) {
		protected.Use(app.IsAuthorized)
		protected.Get("/", app.HelloWorld)
		protected.Get("/logout", app.LogOutUser)
		protected.Get("/deploy", app.Deploy)
		protected.Get("/stop", app.StopInstance)
		protected.Get("/kill", app.DeleteInstance)
		protected.Post("/register", app.RegisterDevice)
		protected.Post("/remove", app.RemoveDevice)
		protected.Post("/createpguser", app.CreatePGUser)
		protected.Post("/createpgdb", app.CreatePgDatabase)
		protected.Post("/deletepguser", app.RemovePgUser)
		protected.Post("/deletepgdb", app.RemovePgDatabase)
	})

	// Public routes (no auth required)
	r.Post("/create_user", app.CreateUser)
	r.Get("/verify", app.VerifyUser)
	r.Post("/login", app.LoginUser)

	return r
}
