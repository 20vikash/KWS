package services

import (
	"context"
	"kws/kws/internal/store"

	"github.com/jackc/pgx/v5"
)

type Services struct {
	PgService interface {
		CreatePostgresUser(ctx context.Context, username, password string) error
		CreateDatabase(ctx context.Context, dbName string, owner string) error
		DropDatabase(ctx context.Context, dbName string) error
		DropPostgresUser(ctx context.Context, uid int, username, password string) error
	}
}

func CreateServices(pgCcon *pgx.Conn, pgMainCon *store.PgServiceStore) *Services {
	return &Services{
		PgService: &PGService{
			Con:    pgCcon,
			MainPg: pgMainCon,
		},
	}
}
