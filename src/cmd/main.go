package main

// TODO: Please do graceful shutdown. Maintain the state.
// Bring all the interfaces down
// Stop all the container and maintain the db state.
// Set IP forward Bit to 0
// Revert all the ip table rules if done any.

// TODO: Startup after boot
// Load all the peer config to the wg kernel module.

import (
	"context"
	"fmt"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	env "kws/kws/internal"
	database "kws/kws/internal/database/connection"
	"kws/kws/internal/docker"
	"kws/kws/internal/docker/services"
	serviceConn "kws/kws/internal/docker/services/connections"
	"kws/kws/internal/mq"
	"kws/kws/internal/store"
	"kws/kws/internal/wg"
	lxd_kws "kws/kws/lxd"
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
	Wg             *wg.WgOperations
	IpAlloc        *wg.IPAllocator
	Services       *services.Services
	LXD            *lxd_kws.LXDKWS
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

	// Initialize mq main instance queue
	queue, err := mq.CreateQueueInstance(mqCh, config.MAIN_INSTANCE_QUEUE, config.RETRY_QUEUE)
	if err != nil {
		log.Fatal("Failed to create instance queue")
	}

	// Initialize mq retry queue
	_, err = mq.CreateRetryQueue(mqCh, config.RETRY_QUEUE)
	if err != nil {
		log.Fatal("Failed to create retry queue")
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

	// Connect to the wireguard server.
	wgCli, err := wg.ConnectToWireguard()
	if err != nil {
		log.Fatal("Cannot connect to the wireguard server.")
	}

	// Create WgOprations struct
	wgOp := &wg.WgOperations{
		Con:        wgCli,
		PrivateKey: env.GetWireguardPrivateKey(),
	}

	// Create IPAllocator
	ipAlloc := &wg.IPAllocator{
		CidrValue:     config.CIDR,
		RedisStore:    &store.RedisStore{Ds: rc},
		WgStore:       &store.WireguardStore{Con: connPool},
		InstanceStore: &store.InstanceStore{Db: connPool},
	}

	docker.IpAlloc = ipAlloc
	docker.Domains = &store.Domain{Con: connPool}

	// Initialize pg service
	pgService := serviceConn.Pg{
		User:     env.GetPGServiceUserName(),
		Password: env.GetPGServicePassword(),
		Host:     env.GetPGServiceHost(),
		Port:     env.GetPGServicePort(),
		Name:     env.GetPGServiceName(),
	}

	// Connect to the pg service
	pgSConn, err := pgService.ConnectToPGServiceBackend(context.Background())
	if err != nil {
		log.Fatal("Failed to connect to pg backend service")
	}

	// Create services instance
	services := services.CreateServices(pgSConn, &store.PgServiceStore{Con: connPool})

	// Connect to LXD
	c, err := lxd_kws.ConnectToLXD()
	if err != nil {
		log.Fatal(err)
	}

	// Create LXDKWS struct instance
	lxdKws := &lxd_kws.LXDKWS{
		Conn: c,
	}

	// Initialize Application
	app := Application{
		Port:           ":8080",
		Store:          store.NewStore(connPool, rc, mqType),
		SessionManager: sessionManager,
		Docker:         docker,
		Mq:             mqType,
		Wg:             wgOp,
		IpAlloc:        ipAlloc,
		Services:       services,
		LXD:            lxdKws,
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

	// Create the main wg0 interface
	err = app.Wg.CreateInterfaceWgMain()
	if err != nil {
		if err.Error() != status.INTERFACE_ALREADY_EXISTS {
			log.Fatal("Cannot create interface wg0")
		}
	}

	// Configure wireguard interface to the kernel module.
	err = app.Wg.ConfigureWireguard()
	if err != nil {
		log.Fatal("Cannot configure wireguard to the kernel module")
	}

	// Set the IP forward bit to 1
	err = app.Wg.SetForwardBitToOne()
	if err != nil {
		log.Fatal("Cannot set forward bit to 1")
	}

	// Start the rabbitmq consumer to listen in the background
	app.ConsumeMessageInstance(app.Mq)

	// HTTP server
	http.ListenAndServe(app.Port, NewRouter(&app))
}
