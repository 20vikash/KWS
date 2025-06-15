package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WireguardStore struct {
	Con *pgxpool.Pool
}

func (wg *WireguardStore) AddPeer(ctx context.Context) error {
	return nil
}

func (wg *WireguardStore) RemovePeer(ctx context.Context) error {
	return nil
}
