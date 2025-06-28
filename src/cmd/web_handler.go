package main

import (
	"html/template"
	"net/http"
)

type Device struct {
	ID        string
	PublicKey string
	IP        string
	Active    bool
}

var templates = template.Must(template.ParseGlob("../web/*.html"))

func (app *Application) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	isAuthorized := app.SessionManager.GetBool(r.Context(), "isAuthorized")
	if isAuthorized {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) RenderSignInPage(w http.ResponseWriter, r *http.Request) {
	isAuthorized := app.SessionManager.GetBool(r.Context(), "isAuthorized")
	if isAuthorized {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := templates.ExecuteTemplate(w, "signin", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) HomeHandler(w http.ResponseWriter, r *http.Request) {
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	data := struct {
		Username string
	}{
		Username: userName,
	}

	err := templates.ExecuteTemplate(w, "home", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) RenderDevicesPage(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Username string
		Devices  []Device // Use your actual Device model type
	}{
		Username: app.SessionManager.GetString(r.Context(), "user_name"),
		Devices:  []Device{{ID: "a", PublicKey: "sassas", IP: "127.0.0.1", Active: true}},
	}

	err := templates.ExecuteTemplate(w, "devices", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
