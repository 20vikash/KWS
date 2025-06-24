package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load .env file variables
func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Cannot load .env variables into the OS")
	}
}

// Redis
func GetRedisHost() string {
	return os.Getenv("REDIS_HOST")
}

func GetRedisPort() string {
	return os.Getenv("REDIS_PORT")
}

func GetRedisPassword() string {
	return os.Getenv("REDIS_PASSWORD")
}
