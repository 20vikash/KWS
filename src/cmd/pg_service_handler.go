package main

import (
	"encoding/json"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/models"
	"log"
	"net/http"
)

type ResponseCreateUser struct {
	ID          int
	Username    string
	Password    string
	Permissions string
	UserLimit   int
}

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
	id, err := app.Store.PgService.AddUser(r.Context(), pgUser)
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
		// Revert db state back.
		app.Store.PgService.RemoveUser(r.Context(), models.CreatePgServiceUser(uid, userName, password))
		return
	}

	response := ResponseCreateUser{
		ID:          id,
		Username:    userName,
		Password:    password,
		Permissions: "Limited",
		UserLimit:   config.MAX_SERVICE_DB_USERS,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *Application) CreatePgDatabase(w http.ResponseWriter, r *http.Request) {
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
	dbName := r.FormValue("db_name")

	// Create PGServiceUser struct
	pgUser := models.CreatePgServiceUser(uid, userName, password)

	// Update the main db
	err = app.Store.PgService.AddDatabase(r.Context(), pgUser, &models.PGServiceDatabase{DbName: dbName})
	if err != nil {
		if err.Error() == status.PG_MAX_DB_LIMIT {
			http.Error(w, "db limit exceeded", http.StatusBadRequest)
			return
		}
		if err.Error() == status.PG_DB_ALREDAY_EXISTS {
			http.Error(w, "DB already exists for the user", http.StatusConflict)
			return
		}

		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// Update the service db
	err = app.Services.PgService.CreateDatabase(r.Context(), dbName, userName)
	if err != nil {
		http.Error(w, "failed to create database", http.StatusInternalServerError)
		// Revert db state back.
		app.Store.PgService.RemoveDatabase(r.Context(), models.CreatePgServiceUser(uid, userName, password), &models.PGServiceDatabase{DbName: dbName})
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Successfully created database"))
}

func (app *Application) RemovePgUser(w http.ResponseWriter, r *http.Request) {
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

	// Update the service db
	err = app.Services.PgService.DropPostgresUser(r.Context(), uid, userName, password)
	if err != nil {
		http.Error(w, "failed to remove user", http.StatusInternalServerError)
		return
	}

	// Update the main db
	err = app.Store.PgService.RemoveUser(r.Context(), models.CreatePgServiceUser(uid, userName, password))
	if err != nil {
		http.Error(w, "failed to remove user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully removed user"))
}

func (app *Application) RemovePgDatabase(w http.ResponseWriter, r *http.Request) {
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
	dbName := r.FormValue("db_name")

	// Create PGServiceUser struct
	pgUser := models.CreatePgServiceUser(uid, userName, password)

	// Update the service database
	err = app.Services.PgService.DropDatabase(r.Context(), dbName)
	if err != nil {
		http.Error(w, "failed to remove database", http.StatusInternalServerError)
		return
	}

	// Update the main DB
	err = app.Store.PgService.RemoveDatabase(r.Context(), pgUser, &models.PGServiceDatabase{DbName: dbName})
	if err != nil {
		http.Error(w, "failed to remove pg database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully removed database"))
}
