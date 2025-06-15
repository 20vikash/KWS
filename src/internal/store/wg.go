package store

import (
	"context"
	"kws/kws/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WireguardStore struct {
	Con *pgxpool.Pool
}

func (wg *WireguardStore) AddPeer(ctx context.Context, uid string, wgType *models.WireguardType) error {
	sql := `
		INSERT INTO wgpeer (user_id, public_key, ip_address) VALUES ($1, $2, $3)
	`

	_, err := wg.Con.Exec(ctx, sql,
		uid,
		wgType.PublicKey,
		wgType.IpAddress,
	)
	if err != nil {
		log.Println("Cannot insert wgpeer record")
		return err
	}

	return nil
}

func (wg *WireguardStore) RemovePeer(ctx context.Context, uid string) error {
	return nil
}
