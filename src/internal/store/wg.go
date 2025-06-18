package store

import (
	"context"
	"errors"
	"kws/kws/consts/config"
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

func (wg *WireguardStore) AddPeer(ctx context.Context, uid int, wgType *models.WireguardType) error {
	var numberOfDevices int

	sql := `
		SELECT COUNT(user_id) FROM wgpeer WHERE user_id = $1
	`

	err := wg.Con.QueryRow(ctx, sql, uid).Scan(&numberOfDevices)
	if err != nil {
		log.Println("Cannot find number of users")
		return err
	}

	if numberOfDevices == config.MAX_WG_DEVICES_PER_USER {
		log.Println("Hit the max device count. Could not add more")
		return errors.New(status.WG_DEVICE_LIMIT)
	}

	sql = `
		INSERT INTO wgpeer (user_id, public_key, ip_address) VALUES ($1, $2, $3)
	`

	_, err = wg.Con.Exec(ctx, sql,
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

func (wg *WireguardStore) RemovePeer(ctx context.Context, pubKey string, uid int) (int, error) {
	var ipAddress int

	sql := `
		DELETE FROM wgpeer WHERE public_key = $1 AND user_id = $2 RETURNING ip_address
	`

	err := wg.Con.QueryRow(ctx, sql,
		pubKey,
		uid,
	).Scan(&ipAddress)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("No rows found to delete")
			return -1, errors.New(status.PEER_DOES_NOT_EXIST)
		}
		log.Println("Cannot delete wgpeer record")
		return -1, err
	}

	return ipAddress, nil
}

func (wg *WireguardStore) AllocateNextFreeIP(ctx context.Context, maxHostNumber int, uid int, wgType *models.WireguardType) (int, error) {
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
			return -1, err
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
				if ip > maxHostNumber {
					return errors.New(status.HOST_EXHAUSTION)
				}
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

		if err.Error() == status.HOST_EXHAUSTION {
			log.Println("Cannot allocate IP more than the host portion size")
			return -1, err
		}

		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == "40001" {
			log.Printf("Serialization conflict, retrying... Attempt: %d\n", i+1)
			continue
		}

		log.Println("Transaction failed. Not serializable error")
		return -1, err
	}

	return ip, nil
}
