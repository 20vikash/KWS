package main

import (
	"kws/kws/internal/gmail"
	"kws/kws/models"
	"kws/kws/status"
	"log"
	"net/http"
	"strings"
)

func (app *Application) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

// Endpoint dedicated to web forms.
func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {
	// If already authorized, redirect to home page.
	isAuthorized := app.sessionManager.GetBool(r.Context(), "isAuthorized")
	if isAuthorized {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

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
		return
	}

	// Generate token and set it to the redis store.
	token := app.Store.Auth.GenerateToken(r.Context(), user.Email)
	app.Store.InMemory.SetEmailToken(r.Context(), user.Email, token)

	// Send the gmail to the address.
	go func(email, token string) {
		err := gmail.SendMail(user.Email, token)
		if err != nil {
			log.Println(err.Error())
		}
	}(user.Email, token)
}

// Verify the email token and verify the user.
func (a *Application) VerifyUser(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	value := a.Store.InMemory.GetEmailFromToken(r.Context(), token)

	email := strings.Split(value, ":")[2]

	err := a.Store.InMemory.DeleteEmailToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusGone)
		w.Write([]byte("The link got expied. Try again"))
		return
	}

	err = a.Store.Auth.VerifyUser(r.Context(), email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to verify. Try again"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully verified the email"))
}

func (app *Application) LoginUser(w http.ResponseWriter, r *http.Request) {
	// If already authorized, redirect to home page.
	isAuthorized := app.sessionManager.GetBool(r.Context(), "isAuthorized")
	if isAuthorized {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Parse form fields.
	err := r.ParseForm()
	if err != nil {
		log.Println("There is an error parsing the form")
	}

	// Usermodel instance
	var userModel *models.User

	// Get the field values
	userName := r.FormValue("user_name")
	password := r.FormValue("password")

	// Login the user
	userModel, err = app.Store.Auth.LoginUser(r.Context(), userName, password)
	if err != nil {
		message := ""

		if err.Error() == status.USER_NAME_INVALID || err.Error() == status.WRONG_CREDENTIALS {
			message = "Password or user name is wrong"
		} else if err.Error() == status.USER_NOT_VERIFIED {
			message = "You are not verified."
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))

		return // I forgot you. You fucking trouble maker
	}

	// Put the userID, userName, and isAuthorized into the session store making them authorized
	app.sessionManager.Put(r.Context(), "id", userModel.Id)
	app.sessionManager.Put(r.Context(), "user_name", userModel.User_name)
	app.sessionManager.Put(r.Context(), "isAuthorized", true)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully"))
}

func (app *Application) LogOutUser(w http.ResponseWriter, r *http.Request) {
	// Check if the user is not authorized
	isAuthorized := app.sessionManager.GetBool(r.Context(), "isAuthorized")
	if !isAuthorized {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("You are not authorized to logout"))

		return
	}

	// Destroy the session
	err := app.sessionManager.Destroy(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to logout"))

		return
	}

	// Redirect the user to the login page.
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
