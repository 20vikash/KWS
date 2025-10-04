package store

import (
	"context"
	"errors"
	"kws/kws/consts/config"
	"kws/kws/models"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TunnelStore struct {
	db *pgxpool.Pool
}

func (ts *TunnelStore) CreateTunnel(ctx context.Context, tunnel models.Tunnels) error {
	sql := `
		INSERT INTO tunnels (user_id, domain, is_custom, tunnel_name) VALUES($1, $2, $3, $4)
	`

	_, err := ts.db.Exec(ctx, sql,
		tunnel.UID,
		tunnel.Domain,
		tunnel.IsCustom,
		tunnel.Name,
	)
	if err != nil {
		if strings.Contains(err.Error(), "23505") { // 23505 is the error code for duplicate violation
			log.Println("Error: Duplicate entry (unique constraint violation)")
			return errors.New("Tunnel entry already exists")
		} else {
			log.Println("Cannot create a new Tunnel. Failed at inserting the tunnel details into the table.")
		}

		return err
	}

	return nil
}

func (ts *TunnelStore) DestroyTunnel(ctx context.Context, tunnel models.Tunnels) error {
	sql := `
		DELETE FROM tunnels WHERE tunnel_name = $1 AND user_id = $2
	`

	_, err := ts.db.Exec(ctx, sql,
		tunnel.Name,
		tunnel.UID,
	)
	if err != nil {
		log.Println("Cannot delete tunnel entry")
	}

	return nil
}

func (ts *TunnelStore) GetDomainFromTunnel(ctx context.Context, tunnel models.Tunnels) (string, error) {
	var domainName string

	sql := `
		SELECT domain FROM tunnels WHERE tunnel_name = $1 AND user_id = $2
	`

	err := ts.db.QueryRow(ctx, sql, tunnel.Name, tunnel.UID).Scan(&domainName)
	if err != nil {
		log.Println("Cannot get domain from tunnel")
		return "", errors.New(config.NO_DOMAIN_FOR_TUNNEL)
	}

	return domainName, nil
}
