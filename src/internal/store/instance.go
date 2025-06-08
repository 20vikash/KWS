package store

import "github.com/jackc/pgx/v5/pgxpool"

type InatanceStore struct {
	db *pgxpool.Pool
}
