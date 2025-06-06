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

// Postgres
func GetDBUserName() string {
	return os.Getenv("DB_USERNAME")
}

func GetDBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func GetDBHost() string {
	return os.Getenv("DB_HOST")
}

func GetDBPort() string {
	return os.Getenv("DB_PORT")
}

func GetDBName() string {
	return os.Getenv("DB_DBNAME")
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

// Gmail
func GetGmailAppPassword() string {
	return os.Getenv("GMAIL_PASSWORD")
}

func GetGmail() string {
	return os.Getenv("GMAIL_ADDRESS")
}
