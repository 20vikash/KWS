package env

import "os"

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
