package main

import (
	"encoding/json"
	"kws/kws/consts/config"
	"kws/kws/internal/store"
	"net/http"
)

type InstanceResponse struct {
	JobID  string
	Action string
}

func (app *Application) handleInstanceAction(w http.ResponseWriter, r *http.Request, action string) {
	// Get the session values (uid and username)
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	var insUser string
	var insPassword string

	// Collect username and password if its deploy
	if action == config.DEPLOY {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		insUser = r.FormValue("insUser")
		insPassword = r.FormValue("insPassword")
	}

	// Generate a job ID
	jid := generateHashedJobID(uid, userName)

	// Push the message to the queue.
	err := app.Store.MessageQueue.PushMessageInstance(r.Context(), &store.QueueMessage{
		InsUser:     insUser,
		InsPassword: insPassword,
		UserID:      uid,
		UserName:    userName,
		JobID:       jid,
		Action:      action,
	})
	if err != nil {
		http.Error(w, "failed to handle your request", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&InstanceResponse{
		JobID:  jid,
		Action: action,
	})
}

func (app *Application) Deploy(w http.ResponseWriter, r *http.Request) {
	app.handleInstanceAction(w, r, config.DEPLOY)
}

func (app *Application) StopInstance(w http.ResponseWriter, r *http.Request) {
	app.handleInstanceAction(w, r, config.STOP)
}

func (app *Application) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	app.handleInstanceAction(w, r, config.KILL)
}
