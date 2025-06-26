package store

import (
	"context"
	"errors"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/models"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgServiceStore struct {
	Con *pgxpool.Pool
}

func (pg *PgServiceStore) findPID(ctx context.Context, userName, password string) (int, error) {
	var pid int

	// Extract the ID from the pg username and pg password
	sql := `
		SELECT id FROM pg_service_user WHERE pg_user_name = $1 AND pg_user_password = $2
	`

	err := pg.Con.QueryRow(ctx, sql,
		userName,
		password,
	).Scan(&pid)
	if err != nil {
		log.Println("Cannot find the id from the given pg username and pg password")
		return -1, err
	}

	return pid, nil
}

func (pg *PgServiceStore) AddUser(ctx context.Context, pgUser *models.PGServiceUser) error {
	var count int

	// Check if the user count limit exceeds
	sql := `
		SELECT COUNT(pg_user_name) FROM pg_service_user WHERE user_id = $1
	`

	err := pg.Con.QueryRow(ctx, sql, pgUser.Uid).Scan(&count)
	if err != nil {
		log.Println("Cannot find the number of pg users")
		return err
	}

	if count >= config.MAX_SERVICE_DB_USERS {
		log.Println("Exceeded the pg user limit")
		return errors.New(status.PG_MAX_USER_LIMIT)
	}

	// Insert db record
	sql = `
		INSERT INTO pg_service_user (user_id, pg_user_name, pg_user_password) VALUES ($1, $2, $3)
	`

	_, err = pg.Con.Exec(ctx, sql, pgUser.Uid, pgUser.UserName, pgUser.Password)
	if err != nil {
		// Check if the username already exists.
		if strings.Contains(err.Error(), "23505") {
			log.Println("PG User already exists")
			return errors.New(status.PG_USER_ALREDAY_EXISTS)
		}

		log.Println("Cannot insert pg user data")
		return err
	}

	log.Println("Successfully created pg user")

	return nil
}

func (pg *PgServiceStore) AddDatabase(ctx context.Context, pgUser *models.PGServiceUser, pgDatabase *models.PGServiceDatabase) error {
	pid, err := pg.findPID(ctx, pgUser.UserName, pgUser.Password)
	if err != nil {
		return err
	}

	var dbCount int

	// Check if the Database count limit has exceeded
	sql := `
		SELECT COUNT(id) FROM pg_service_db WHERE pid = $1
	`

	err = pg.Con.QueryRow(ctx, sql, pid).Scan(&dbCount)
	if err != nil {
		log.Println("Cannot find the count of databases for the pg user")
		return err
	}

	if dbCount >= config.MAX_SERVICE_DB_DB {
		log.Println("DB limit for user exceeded")
		return errors.New(status.PG_MAX_DB_LIMIT)
	}

	// Update the db record
	sql = `
		INSERT INTO pg_service_db (pid, db_name) VALUES($1, $2)
	`

	_, err = pg.Con.Exec(ctx, sql, pid, pgDatabase.DbName)
	if err != nil {
		// Check if the database already exists.
		if strings.Contains(err.Error(), "23505") {
			log.Println("PG DB already exists")
			return errors.New(status.PG_DB_ALREDAY_EXISTS)
		}

		return err
	}

	log.Println("Successfully created pg service database")

	return nil
}

func (pg *PgServiceStore) RemoveDatabase(ctx context.Context, pgUser *models.PGServiceUser, pgDatabase *models.PGServiceDatabase) error {
	var pid int

	// Verify user identity
	sql := `
		SELECT id FROM pg_service_user
		WHERE pg_user_name = $1 AND pg_user_password = $2 AND user_id = $3
	`

	err := pg.Con.QueryRow(ctx, sql, pgUser.UserName, pgUser.Password, pgUser.Uid).Scan(&pid)
	if err != nil {
		log.Println("Could not find pg user for deletion")
		return errors.New(status.PG_USER_NOT_FOUND)
	}

	// Delete the database record
	sql = `
		DELETE FROM pg_service_db WHERE pid = $1 AND db_name = $2
	`

	res, err := pg.Con.Exec(ctx, sql, pid, pgDatabase.DbName)
	if err != nil {
		log.Printf("Failed to delete database %s for pid=%d: %v\n", pgDatabase.DbName, pid, err)
		return err
	}

	if res.RowsAffected() == 0 {
		log.Println("No matching database found to delete")
		return errors.New(status.PG_DB_NOT_FOUND)
	}

	log.Printf("Deleted database '%s' for pg user '%s'\n", pgDatabase.DbName, pgUser.UserName)
	return nil
}

func (pg *PgServiceStore) RemoveUser(ctx context.Context, pgUser *models.PGServiceUser) error {
	// Delete the user record
	sql := `
		DELETE FROM pg_service_user
		WHERE pg_user_name = $1 AND pg_user_password = $2 AND user_id = $3
	`

	res, err := pg.Con.Exec(ctx, sql, pgUser.UserName, pgUser.Password, pgUser.Uid)
	if err != nil {
		log.Printf("Failed to delete pg user %s: %v\n", pgUser.UserName, err)
		return err
	}

	if res.RowsAffected() == 0 {
		log.Printf("No matching user '%s' found for deletion\n", pgUser.UserName)
		return errors.New(status.PG_USER_NOT_FOUND)
	}

	log.Printf("Deleted PG user '%s' and their databases\n", pgUser.UserName)
	return nil
}

func (pg *PgServiceStore) GetUserDatabases(ctx context.Context, userName, password string) ([]string, error) {
	var dbName string
	var dbNames = make([]string, 0)

	pid, err := pg.findPID(ctx, userName, password)
	if err != nil {
		return nil, err
	}

	sql := `
		SELECT db_name FROM pg_service_db WHERE pid = $1
	`

	rows, err := pg.Con.Query(ctx, sql, pid)
	if err != nil {
		log.Println("Cannot get all the databases of the user")
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&dbName); err != nil {
			return nil, err
		}

		dbNames = append(dbNames, dbName)
	}

	return dbNames, nil
}
