package store

import (
	"context"
	"kws/kws/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InatanceStore struct {
	db *pgxpool.Pool
}

func (i *InatanceStore) CreateInstance(ctx context.Context, uid int, userName string) error {
	instance := models.CreateInstanceType(uid, userName)

	sql := `
		INSERT INTO instance (user_id, volume_name, container_name, instance_type, is_running) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := i.db.Exec(ctx, sql,
		instance.Uid,
		instance.VolumeName,
		instance.ContainerName,
		instance.InstanceType,
		true,
	)
	if err != nil {
		log.Println("Cannot insert row into instance table")
		return err
	}

	return nil
}
