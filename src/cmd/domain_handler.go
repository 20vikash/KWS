package main

import (
	"encoding/json"
	"kws/kws/consts/config"
	"kws/kws/internal/store"
	"kws/kws/models/web"
	"log"
	"net/http"
	"strconv"
)

type DomainResponse struct {
	JobID  string `json:"jobID"`
	Action string `json:"action"`
}

func (app *Application) handleDomainAction(w http.ResponseWriter, r *http.Request, action string) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Failed to parse form")
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	domain := r.FormValue("domain_name")

	var port int
	if action == config.ADD_USER_DOMAIN {
		portStr := r.FormValue("port")

		port, err = strconv.Atoi(portStr)
		if err != nil {
			http.Error(w, "Invalid port number", http.StatusBadRequest)
			return
		}
	}

	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	// Generate a job ID
	jid := generateHashedJobID(uid, userName)

	// Push the message to the queue.
	err = app.Store.MessageQueue.PushMessageInstance(r.Context(), &store.DomainQueueMessage{

		JobID:  jid,
		Domain: domain,
		Port:   port,
		UserID: uid,
		Action: action,
	},
		app.MqPool,
	)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&DomainResponse{
		JobID:  jid,
		Action: action,
	})

}
func (app *Application) AddUserDomain(w http.ResponseWriter, r *http.Request) {
	app.handleDomainAction(w, r, config.ADD_USER_DOMAIN)
}

func (app *Application) RemoveUserDomain(w http.ResponseWriter, r *http.Request) {
	app.handleDomainAction(w, r, config.REMOVE_USER_DOMAIN)
}

func (app *Application) DomainResult(w http.ResponseWriter, r *http.Request) {
	app.handleDomainResult(w, r)
}

func (app *Application) handleDomainResult(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobID")

	done, result, err := app.Store.InMemory.GetUserDomainResult(r.Context(), jobID)
	if err != nil {
		http.Error(w, "failed to handle your request", http.StatusInternalServerError)
		return
	}

	if result == nil {
		result = &web.JobResponseDomain{}
	}

	result.Done = done

	writeJSON(w, result)
}
