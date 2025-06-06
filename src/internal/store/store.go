package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Auth interface {
		CreateUser(ctx context.Context, first_name, last_name, email, password string) error
	}
}

func NewStore(pg *pgxpool.Pool) *Storage {
	return &Storage{
		Auth: &AuthStore{
			db: pg,
		},
	}
}
