package services

import (
	"github.com/jackc/pgx/v5"
)

type Services struct {
	PgService interface {
	}
}

func CreateServices(con *pgx.Conn) *Services {
	return &Services{
		PgService: PGService{Con: con},
	}
}
