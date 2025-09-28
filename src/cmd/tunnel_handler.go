package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"kws/kws/consts/status"
	"net/http"
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
