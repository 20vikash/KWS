package main

import (
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/models"
	"net/http"
)

func (app *Application) Deploy(w http.ResponseWriter, r *http.Request) {
	//TODO: Pass this request to a message queue, and process it all in the background. SSE all the info.

	// This is just a test code. Should change this.
	uid := app.SessionManager.GetInt(r.Context(), "id")
	userName := app.SessionManager.GetString(r.Context(), "user_name")

	instanceType := models.CreateInstanceType(uid, userName)
	id, err := app.Docker.CreateContainerCore(r.Context(),
		instanceType.ContainerName,
		instanceType.VolumeName,
		config.CORE_NETWORK_NAME,
	)
	if err != nil {
		http.Error(w, "cannot deploy instance", http.StatusInternalServerError)
		return
	}

	err = app.Docker.StartContainer(r.Context(), id)
	if err != nil {
		if err.Error() == status.CONTAINER_ALREADY_RUNNING {
			http.Error(w, "instance is already running", http.StatusBadRequest)
			return
		}

		http.Error(w, "cannot start instance", http.StatusInternalServerError)
		return
	}
}

func (app *Application) StopInstance(w http.ResponseWriter, r *http.Request) {
	//TODO: Same as above
}

func (app *Application) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	//TODO: Not visible in the UI, CRON jobs will use this endpoint
}
