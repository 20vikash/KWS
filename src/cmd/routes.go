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
	// r.Use(app.LoginRateLimitMiddleware)

	// Define a sub-router for protected tunnel routes
	r.Group(func(protected chi.Router) {
		protected.Use(app.IsTunnelUserAuthorized)
		protected.Post("/create_tunnel", app.CreateTunnel)
		protected.Post("/destroy_tunnel", app.DestroyTunnel)
	})

	// Define a sub-router for protected routes
	r.Group(func(protected chi.Router) {
		protected.Use(app.IsAuthorized)
		protected.Get("/logout", app.LogOutUser)
		protected.Post("/deploy", app.Deploy)
		protected.Get("/stop", app.StopInstance)
		protected.Get("/kill", app.DeleteInstance)
		protected.Post("/register", app.RegisterDevice)
		protected.Post("/remove", app.RemoveDevice)
		protected.Post("/createpguser", app.CreatePGUser)
		protected.Post("/createpgdb", app.CreatePgDatabase)
		protected.Post("/deletepguser", app.RemovePgUser)
		protected.Post("/deletepgdb", app.RemovePgDatabase)
		protected.Get("/kws_devices", app.RenderDevicesPage)
		protected.Post("/active", app.IsOnline)
		protected.Post("/deployresult", app.DeployResult)
		protected.Post("/stopresult", app.StopResult)
		protected.Post("/killresut", app.KillResult)
		protected.Post("/adddomain", app.AddUserDomain)
		protected.Post("/removedomain", app.RemoveUserDomain)
		protected.Route("/kws_services", func(r chi.Router) {
			r.Get("/", app.RenderServicesPage)
			r.Get("/postgres/users", app.RenderPgUsersPage)
			r.Get("/postgres/db", app.RenderPgDatabasesPage)
		})
		protected.Get("/kws_instances", app.RenderInstancePage)
		protected.Get("/kws_publish", app.RenderPublishPage)
		protected.Get("/", app.HomeHandler)
	})

	// Public routes (no auth required)
	r.Post("/create_user", app.CreateUser)
	r.Get("/verify", app.VerifyUser)
	r.Post("/login", app.LoginUser)
	r.Get("/kws_register", app.RenderRegisterPage)
	r.Get("/kws_signin", app.RenderSignInPage)
	r.Post("/tunnel_login", app.LoginTunnelUser)

	// Serve static files
	r.Handle("/js/*", http.StripPrefix("/js/", http.FileServer(http.Dir("../web/js"))))
	r.Handle("/css/*", http.StripPrefix("/css/", http.FileServer(http.Dir("../web/css"))))

	return r
}
