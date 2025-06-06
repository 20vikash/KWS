package env

import "os"

func GetDBUserName() string {
	return os.Getenv("DB_USERNAME")
}

// Postgres
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
