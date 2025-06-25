package store

import (
	"context"
	"errors"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"kws/kws/models"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgServiceStore struct {
	Con *pgxpool.Pool
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

	if count == config.MAX_SERVICE_DB_USERS {
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
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				log.Println("User already exists")
				return errors.New(status.PG_USER_ALREDAY_EXISTS)
			}
		}

		log.Println("Cannot insert pg user data")
		return err
	}

	return nil
}

func (pg *PgServiceStore) AddDatabase(ctx context.Context, pgDatabase *models.PGServiceDatabase) error {
	return nil
}
