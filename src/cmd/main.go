package main

import (
	"context"
	"fmt"
	env "kws/kws/internal"
	database "kws/kws/internal/database/connection"
	"kws/kws/internal/docker"
	"kws/kws/internal/mq"
	"kws/kws/internal/store"
	"log"
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
	SessionManager *scs.SessionManager
	Docker         *docker.Docker
	Mq             *store.MQ
}

func main() {
	// Load .env variables into OS.
	env.LoadEnv()

	// Get dockerCon connection
	dockerCon, err := docker.GetConnection()
	if err != nil {
		log.Fatal("Failed to connect to docker")
	}
	docker := &docker.Docker{
		Con: dockerCon,
	}

	// Get rabbitmq connection and set up channel.
	mq := mq.Mq{
		User: env.GetMqUser(),
		Pass: env.GetMqPassword(),
		Port: env.GetMqPort(),
		Host: env.GetMqHost(),
	}
	con, err := mq.ConnectToMq() // TCP connection
	if err != nil {
		log.Fatal("Failed to connect to rabbitmq")
	}
	mqCh, err := mq.CreateChannel(con) // Channel connection
	if err != nil {
		log.Fatal("Failed to create a Mq channel")
	}

	// Initialize mq queue
	queue, err := mq.CreateQueue(mqCh, "instance")
	if err != nil {
		log.Fatal("Failed to create instance queue")
	}

	// Create a consumer for that queue
	consumer, err := mq.CreateConsumer(mqCh, queue)
	if err != nil {
		log.Fatal("Failed to create a consumer")
	}

	// Create MQ struct instance.
	mqType := &store.MQ{
		Consumer: consumer,
		Ch:       mqCh,
		Queue:    queue,
	}

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
		Store:          store.NewStore(connPool, rc, mqType),
		SessionManager: sessionManager,
		Docker:         docker,
	}

	// Initialize the server with the docker images
	err = docker.CreateImageCore(context.Background())
	if err != nil {
		log.Fatal("Core image creation error")
	}

	// Initialize the server with custom bridge networks
	err = docker.CreateCustomNetwork(context.Background())
	if err != nil {
		log.Fatal("Core network creation error")
	}

	// Start the rabbitmq consumer to listen in the background
	app.ConsumeMessageInstance(app.Mq)

	// HTTP server
	http.ListenAndServe(app.Port, NewRouter(&app))
}
