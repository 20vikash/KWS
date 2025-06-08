package store

import (
	"context"
	"kws/kws/models"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InstanceStore struct {
	db *pgxpool.Pool
}

// Database level function: Insert an instance record related to the user
func (i *InstanceStore) CreateInstance(ctx context.Context, uid int, userName string) error {
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

// Database level function: Remove an instance record related to the user
func (i *InstanceStore) RemoveInstance(ctx context.Context, uid int) error {
	sql := `
		DELETE FROM instance WHERE user_id = $1
	`

	_, err := i.db.Exec(ctx, sql, uid)
	if err != nil {
		log.Println("Cannot stop/remove the instance")
		return err
	}

	return nil
}
