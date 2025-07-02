package store

import (
	"context"
	"kws/kws/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Domain struct {
	Con *pgxpool.Pool
}

func (d *Domain) AddDomain(ctx context.Context, domain *models.Domain) error {
	return nil
}
