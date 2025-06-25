package store

import (
	"kws/kws/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgServiceStore struct {
	Con *pgxpool.Pool
}

func (pg *PgServiceStore) AddUser(pgUser *models.PGServiceUser) error {
	return nil
}

func (pg *PgServiceStore) AddDatabase(pgDatabase *models.PGServiceDatabase) error {
	return nil
}
