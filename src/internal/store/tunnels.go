package store

import (
	"context"
	"errors"
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
		INSERT INTO tunnels (user_id, domain, is_custom) VALUES($1, $2, $3)
	`

	_, err := ts.db.Exec(ctx, sql,
		tunnel.UID,
		tunnel.Domain,
		tunnel.IsCustom,
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
		DELETE FROM tunnels WHERE uid = $1
	`

	_, err := ts.db.Exec(ctx, sql,
		tunnel.UID,
	)
	if err != nil {
		log.Println("Cannot delete tunnel entry")
	}

	return nil
}
