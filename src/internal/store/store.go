package store

import "github.com/jackc/pgx/v5/pgxpool"

type Storage struct {
	Auth interface {
		Test()
	}
}

func NewStore(pg *pgxpool.Pool) *Storage {
	return &Storage{
		Auth: &AuthStore{
			db: pg,
		},
	}
}
