package main

import (
	"fmt"
	"net/http"
	"proxy/proxy/env"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

var sessionManager *scs.SessionManager

type Application struct {
	SessionManager *scs.SessionManager
	Port           string
}

func main() {
	// Load .env variables into OS.
	env.LoadEnv()

	// Set up redis db pool for session manager.
	rPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				fmt.Sprintf("%s:%s",
					env.GetRedisHost(), env.GetRedisPort()),
				redis.DialDatabase(1),
				redis.DialPassword(env.GetRedisPassword()),
			)
		},
	}

	// Initialize session manager.
	sessionManager = scs.New()
	sessionManager.Store = redisstore.New(rPool)

	// Initialize application
	app := &Application{
		SessionManager: sessionManager,
		Port:           ":9000",
	}

	http.ListenAndServe(app.Port, NewRouter(app))
}
