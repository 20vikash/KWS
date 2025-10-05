package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/internal/nginx"
	"kws/kws/models"
	"log"
	"net/http"
	"strconv"
)

func (app *Application) LoginTunnelUser(w http.ResponseWriter, r *http.Request) {
	// Parse form fields
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	// Get the field values
	userName := r.FormValue("user_name")
	password := r.FormValue("password")

	// Login the user
	userModel, err := app.Store.Auth.LoginUser(r.Context(), userName, password)
	if err != nil {
		message := ""
		if err.Error() == status.USER_NAME_INVALID || err.Error() == status.WRONG_CREDENTIALS {
			message = "Password or user name is wrong"
		} else if err.Error() == status.USER_NOT_VERIFIED {
			message = "You are not verified."
		}

		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// Generate a new secret (use crypto/rand for cryptographic randomness)
	secretBytes := make([]byte, 32) // 256-bit secret
	if _, err := rand.Read(secretBytes); err != nil {
		http.Error(w, "Failed to generate secret", http.StatusInternalServerError)
		return
	}
	secret := hex.EncodeToString(secretBytes)

	// Save secret in Redis with the userId
	err = app.Store.InMemory.SetTunnelLogin(r.Context(), secret, userModel.Id)
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// Prepare response
	resp := map[string]any{
		"status": "success",
		"secret": secret,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (app *Application) CreateTunnel(w http.ResponseWriter, r *http.Request) {
	// Parse form fields
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	uid := r.Context().Value("uid").(int)

	// Get the field values
	domain := r.FormValue("domain")
	isCustom := r.FormValue("is_custom")
	tunnelName := r.FormValue("name")

	isCustomB, err := strconv.ParseBool(isCustom)
	if err != nil {
		log.Println("Cannot parse isCustom bool")
		http.Error(w, "Bad format", http.StatusBadRequest)
		return
	}

	// TODO: Make it all async using MQ
	// Create nginx conf

	if isCustomB {
		// TODO: Use certbot to generate SSL certs
	}

	template := nginx.Template{
		Domain: domain,
	}
	err = template.AddNewConf(config.DOMAIN_TEMPLATE)
	if err != nil {
		log.Println("Failed to create nginx conf file for tunnel")
		http.Error(w, "cannot create nginx conf file", http.StatusInternalServerError)
		return
	}

	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
	if err != nil {
		template.RemoveConf()
		app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)
		log.Println("Failed to create nginx conf file for tunnel")
		http.Error(w, "cannot create nginx conf file", http.StatusInternalServerError)
		return
	}

	err = app.Store.Tunnels.CreateTunnel(r.Context(), models.Tunnels{
		UID:      uid,
		Domain:   domain,
		IsCustom: isCustomB,
		Name:     tunnelName,
	})

	if err != nil {
		log.Println("Something went wrong with creating a tunnel(handler)")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (app *Application) DestroyTunnel(w http.ResponseWriter, r *http.Request) {
	// Parse form fields
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	uid := r.Context().Value("uid").(int)

	tunnelName := r.FormValue("name")

	err = app.Store.Tunnels.DestroyTunnel(r.Context(), models.Tunnels{
		UID:  uid,
		Name: tunnelName,
	})

	if err != nil {
		log.Println("Cannot delete tunnel (handler)")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Remove nginx conf
	domain, err := app.Store.Tunnels.GetDomainFromTunnel(r.Context(), models.Tunnels{
		Name: tunnelName,
	})
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	template := nginx.Template{
		Domain: domain,
	}

	err = template.RemoveConf()
	err = app.Docker.ReloadNginxConf(config.NGINX_CONTAINER)

	if err != nil {
		log.Println("Something went wrong in removing nginx conf for tunnel")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
}
