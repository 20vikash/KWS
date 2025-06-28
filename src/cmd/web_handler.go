package main

import (
	"html/template"
	"net/http"
)

type Device struct {
	PublicKey string
	IP        string
	Active    bool
}

type Data struct {
	Username string
	Devices  []Device
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
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Fetch peers information
	peers, err := app.Store.Wireguard.GetDevices(r.Context(), uid)
	if err != nil {
		http.Error(w, "cannot load this page right now", http.StatusInternalServerError)
	}

	var data = new(Data)

	for _, peer := range peers {
		ipAddress := app.IpAlloc.GenerateIP(peer.IpAddress)
		data.Devices = append(data.Devices, Device{PublicKey: peer.PublicKey, IP: ipAddress, Active: false})
	}

	data.Username = userName

	err = templates.ExecuteTemplate(w, "devices", data)
	if err != nil {
		http.Error(w, "Template rendering error", http.StatusInternalServerError)
	}
}
