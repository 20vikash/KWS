package main

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseGlob("../web/*.html"))

func (app *Application) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "register", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *Application) RenderSignInPage(w http.ResponseWriter, r *http.Request) {
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
