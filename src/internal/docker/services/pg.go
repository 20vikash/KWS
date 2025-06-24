package services

import (
	"github.com/jackc/pgx/v5"
)

type PGService struct {
	Con *pgx.Conn
}
