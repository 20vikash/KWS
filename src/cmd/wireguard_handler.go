package main

import (
	"encoding/json"
	"net/http"
)

type RequestBody struct {
	PublicKey string `json:"public_key"`
}

func (app *Application) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	var rb RequestBody

	// Read User ID from the session token
	uid := app.SessionManager.GetInt(r.Context(), "id")

	// Decode JSON body into struct to get the public key
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = app.Wg.AddPeer(r.Context(), uid, rb.PublicKey, app.IpAlloc)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully added peer"))
	// TODO: Return a JSON response with the allocated IP
}

func (app *Application) RemoveDevice(w http.ResponseWriter, r *http.Request) {
	var rb RequestBody

	// Read User ID from the session token
	uid := app.SessionManager.GetInt(r.Context(), "id")

	// Decode JSON body into struct to get the public key
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Remove peer
	err = app.Wg.RemovePeer(r.Context(), rb.PublicKey, uid, app.IpAlloc)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully removed a peer"))
}
