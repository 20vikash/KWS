package store

import "github.com/jackc/pgx/v5/pgxpool"

type AuthStore struct {
	db *pgxpool.Pool
}

func (auth *AuthStore) Test() {
	return
}
