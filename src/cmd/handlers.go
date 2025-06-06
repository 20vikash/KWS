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
	if !(user.ValidateEmail() && user.ValidatePassword() && user.ValidateUserName() && user.ValidateFLNames()) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("One of the field values are not in the right format"))
		return
	}

	err = app.Store.Auth.CreateUser(r.Context(),
		user.First_name,
		user.Last_name,
		user.Email,
		user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Println(email)
}
