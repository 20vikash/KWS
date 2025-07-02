package store

import (
	"context"
	"errors"
	"kws/kws/consts/status"
	"kws/kws/models"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Domain struct {
	Con *pgxpool.Pool
}

func (d *Domain) AddDomain(ctx context.Context, domain *models.Domain) error {
	sql := `
		INSERT INTO domains(user_id, domain_name, port)
		VALUES ($1, $2, $3)
	`

	_, err := d.Con.Exec(ctx, sql, domain.Uid, domain.Domain, domain.Port)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			log.Println("Domain already exists")
			return errors.New(status.DOMAIN_ALREADY_EXISTS)
		}
		log.Println("Cannot insert domain data")
		return err
	}

	return nil
}

func (d *Domain) RemoveDomain(ctx context.Context, domain *models.Domain) error {
	sql := `
		DELETE FROM domains WHERE domain_name = $1 AND user_id = $2
	`

	_, err := d.Con.Exec(ctx, sql,
		domain.Domain,
		domain.Uid,
	)
	if err != nil {
		log.Println("Cannot delete domain record")
		return err
	}

	return nil
}
