package serviceConn

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Pg struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func (p *Pg) ConnectToPGServiceBackend(ctx context.Context) (*pgx.Conn, error) {
	url := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.Name)

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		log.Println("Cannot connect to pg service")
		return nil, err
	}

	return conn, nil
}
