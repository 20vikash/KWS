package main

import (
	"kws/kws/consts/status"
	"kws/kws/models"
	"log"
	"net/http"
)

func (app *Application) CreatePGUser(w http.ResponseWriter, r *http.Request) {
	// Get the user ID
	uid := app.SessionManager.GetInt(r.Context(), "id")

	// Parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println("There is an error parsing the form")
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	// Read the form values.
	userName := r.FormValue("user_name")
	password := r.FormValue("password")

	// Create PGServiceUser struct
	pgUser := models.CreatePgServiceUser(uid, userName, password)

	// Update the main DB
	err = app.Store.PgService.AddUser(r.Context(), pgUser)
	if err != nil {
		if err.Error() == status.PG_MAX_USER_LIMIT {
			http.Error(w, "user limit exceeded", http.StatusBadRequest)
			return
		}
		if err.Error() == status.PG_USER_ALREDAY_EXISTS {
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Update the users pg service db
	err = app.Services.PgService.CreatePostgresUser(r.Context(), userName, password)
	if err != nil {
		http.Error(w, "failed to create pg user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully created pg user"))
}

func (app *Application) CreatePgDatabase(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) RemovePgUser(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) RemovePgDatabase(w http.ResponseWriter, r *http.Request) {

}
