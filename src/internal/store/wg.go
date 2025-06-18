package store

import (
	"context"
	"errors"
	"kws/kws/consts/status"
	"kws/kws/models"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (wg *WireguardStore) AllocateNextMaxIP(ctx context.Context, uid string, wgType *models.WireguardType) error {
	var ip int
	maxRetries := 5

	sqlSelect := `
		SELECT ip_address FROM wgpeer ORDER BY ip_address DESC LIMIT 1
	`
	sqlInsert := `
		INSERT INTO wgpeer (user_id, public_key, ip_address) VALUES ($1, $2, $3)
	`

	for i := range maxRetries {
		tx, err := wg.Con.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
		if err != nil {
			log.Println("Cannot start transaction")
			return err
		}

		err = func() error {
			defer tx.Rollback(ctx)

			err := tx.QueryRow(ctx, sqlSelect).Scan(&ip)
			if err != nil {
				if err == pgx.ErrNoRows {
					log.Println("Cannot find the max of the ip. This is the start it seems")
					ip = 2
				} else {
					log.Println("Cannot find max of the IP. Something went wrong.")
					return err
				}
			} else {
				ip += 1
			}

			_, err = tx.Exec(ctx, sqlInsert,
				uid,
				wgType.PublicKey,
				ip,
			)
			if err != nil {
				log.Println("Cannot insert ip+1 record")
				return err
			}

			return tx.Commit(ctx)
		}()

		if err == nil {
			log.Println("Transaction successful")
			break
		}

		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == "40001" {
			log.Printf("Serialization conflict, retrying... Attempt: %d\n", i+1)
			continue
		}

		log.Println("Transaction failed. Not serializable error")
		return err
	}

	return nil
}
