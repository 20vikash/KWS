package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseBody struct {
	Active    bool   `json:"active"`
	Ip        string `json:"ip"`
	PublicKey string `json:"publicKey"`
}

func (app *Application) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println("There is an error parsing the form")
	}

	// Read User ID from the session token
	uid := app.SessionManager.GetInt(r.Context(), "id")

	// Read form value.
	publicKey := r.FormValue("public_key")

	ipAddr, err := app.Wg.AddPeer(r.Context(), uid, publicKey, app.IpAlloc)
	if err != nil {
		http.Error(w, "Cannot add device", http.StatusInternalServerError)
		return
	}

	response := ResponseBody{
		Active:    false,
		Ip:        ipAddr,
		PublicKey: publicKey,
	}

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *Application) RemoveDevice(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println("There is an error parsing the form")
	}

	// Read User ID from the session token
	uid := app.SessionManager.GetInt(r.Context(), "id")

	publicKey := r.FormValue("public_key")

	// Remove peer
	err = app.Wg.RemovePeer(r.Context(), publicKey, uid, app.IpAlloc)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully removed a peer"))
}
