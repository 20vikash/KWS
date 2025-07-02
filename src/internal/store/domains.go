package store

import (
	"context"
	"errors"
	"kws/kws/consts/status"
	"kws/kws/models"
	"kws/kws/models/web"
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

func (d *Domain) GetUserDomains(ctx context.Context, domain *models.Domain) (*[]web.Domain, error) {
	domains := new([]web.Domain)

	var domainName string
	var port int

	sql := `
		SELECT domain_name, port FROM domains WHERE user_id = $1 AND is_code = $2
	`

	rows, err := d.Con.Query(ctx, sql,
		domain.Uid,
		false,
	)
	if err != nil {
		log.Println("Cannot get all the domains of the user")
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&domainName, &port)
		if err != nil {
			log.Println("Error scanning domain table")
			return nil, err
		}

		*domains = append(*domains, web.Domain{Name: domainName, Port: port, Status: "Active"})
	}

	return domains, nil
}
