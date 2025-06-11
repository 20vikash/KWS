package main

import (
	"encoding/json"
	"kws/kws/internal/store"
	"net/http"
)

type InstanceResponse struct {
	JobID string
}

func (app *Application) Deploy(w http.ResponseWriter, r *http.Request) {
	// Get the session values (uid and username)
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Generate a job ID
	jid := generateHashedJobID(uid, userName)

	// Push the message to the queue.
	err := app.Store.MessageQueue.PushMessageInstance(r.Context(), &store.QueueMessage{
		UserID:   uid,
		UserName: userName,
		JobID:    jid,
	})
	if err != nil {
		http.Error(w, "failed to handle your request", http.StatusInternalServerError)
	}

	instanceResponse := &InstanceResponse{
		JobID: jid,
	}

	// Send the json response with the Job ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(instanceResponse)
}

func (app *Application) StopInstance(w http.ResponseWriter, r *http.Request) {
	//TODO: Same as above
}

func (app *Application) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	//TODO: Not visible in the UI, CRON jobs will use this endpoint
}
