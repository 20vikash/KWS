package main

import (
	"kws/kws/models"
	"log"
	"net/http"
)

func (app *Application) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

// Endpoint dedicated to web forms.
func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseForm()
	if err != nil {
		log.Println("There is an error parsing the form")
	}

	// Read the form values
	email := r.FormValue("email")
	password := r.FormValue("password")
	userName := r.FormValue("user_name")
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")

	// Create User struct
	user := models.User{
		Email:      email,
		Password:   password,
		User_name:  userName,
		First_name: firstName,
		Last_name:  lastName,
	}

	// Validate before updating the database.
	if !(user.ValidateEmail()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email is in the wrong format"))
		return
	}
	if !(user.ValidatePassword()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password format is wrong"))
		return
	}
	if !(user.ValidateFLNames()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("First name or/and last name is in the wrong format"))
		return
	}
	if !(user.ValidateUserName()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User name format is wrong"))
		return
	}

	// Update the database after all the validations
	err = app.Store.Auth.CreateUser(r.Context(),
		user.First_name,
		user.Last_name,
		user.Email,
		user.Password,
		user.User_name,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}
