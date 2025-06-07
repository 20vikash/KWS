package main

import (
	"fmt"
	env "kws/kws/internal"
	database "kws/kws/internal/database/connection"
	"kws/kws/internal/store"
	"net/http"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

var sessionManager *scs.SessionManager

type Application struct {
	Port           string
	Store          *store.Storage
	sessionManager *scs.SessionManager
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

	// Session manager cookie properties.
	sessionManager.Lifetime = 24 * time.Hour                  // Session cookie timeout
	sessionManager.Cookie.Name = "kws_session"                // Session cookie name
	sessionManager.Cookie.HttpOnly = true                     // Javascript cannot read the cookie
	sessionManager.Cookie.Persist = true                      // Persists after browser restart
	sessionManager.Cookie.SameSite = http.SameSiteDefaultMode // Only send the session cookie if I am in the same site.
	sessionManager.Cookie.Secure = env.IsProd()               // Set in the .env (HTTPS mode)

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
		Addr:     fmt.Sprintf("%s:%s", env.GetRedisHost(), env.GetRedisPort()),
		Password: env.GetRedisPassword(),
		DB:       0,
	}
	rc := redis.Connect()

	// Initialize Application
	app := Application{
		Port:           ":8080",
		Store:          store.NewStore(connPool, rc),
		sessionManager: sessionManager,
	}

	// HTTP server
	http.ListenAndServe(app.Port, NewRouter(&app))
}
