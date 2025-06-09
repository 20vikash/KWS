package main

import "net/http"

func (app *Application) Deploy(w http.ResponseWriter, r *http.Request) {
	//TODO: Pass this request to a message queue, and process it all in the background. SSE all the info.
}

func (app *Application) StopInstance(w http.ResponseWriter, r *http.Request) {
	//TODO: Same as above
}

func (app *Application) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	//TODO: Not visible in the UI, CRON jobs will use this endpoint
}
