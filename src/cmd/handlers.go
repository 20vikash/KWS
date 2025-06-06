package main

import (
	"log"
	"net/http"
)

func (app *Application) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}

// Endpoint dedicated to web forms.
func (app *Application) CreateUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("There is an error parsing the form")
	}

	email := r.FormValue("email")
	log.Println(email)
}
