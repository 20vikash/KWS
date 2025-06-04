package main

import (
	"net/http"
)

type Application struct {
	Port string
}

func main() {
	// Initialize Application
	app := Application{
		Port: ":8080",
	}

	// HTTP server
	http.ListenAndServe(app.Port, NewRouter(&app))
}
