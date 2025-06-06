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

	// Initialize Redis database
	redis := database.RedisDB{
		Addr:     "redis_db:6379",
		Password: env.GetRedisPassword(),
		DB:       0,
	}
	rc := redis.Connect()

	// Initialize Application
	app := Application{
		Port:  ":8080",
		Store: store.NewStore(connPool, rc),
	}

	// HTTP server
	http.ListenAndServe(app.Port, NewRouter(&app))
}
