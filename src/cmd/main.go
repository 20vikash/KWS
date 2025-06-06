package main

import (
	env "kws/kws/internal"
	database "kws/kws/internal/database/connection"
	"kws/kws/internal/store"
	"net/http"
)

type Application struct {
	Port  string
	Store *store.Storage
}

func main() {
	// Initialize Pg database
	pg := database.Pg{
		User:     env.GetDBUserName(),
		Password: env.GetDBPassword(),
		Host:     env.GetDBHost(),
		Port:     env.GetDBPort(),
		Name:     env.GetDBName(),
	}
	connPool := pg.GetNewDBConnection()

	// Initialize Application
	app := Application{
		Port:  ":8080",
		Store: store.NewStore(connPool),
	}

	// HTTP server
	http.ListenAndServe(app.Port, NewRouter(&app))
}
