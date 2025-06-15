package store

import (
	"context"
	"errors"
	"kws/kws/consts/status"
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
	sql := `
		DELETE FROM wgpeer WHERE user_id = $1
	`

	rows, err := wg.Con.Exec(ctx, sql, uid)
	if err != nil {
		log.Println("Cannot delete wgpeer record")
		return err
	}

	rowsAffected := rows.RowsAffected()
	if rowsAffected == 0 {
		log.Println("No rows found to delete")
		return errors.New(status.PEER_DOES_NOT_EXIST)
	}

	return nil
}
