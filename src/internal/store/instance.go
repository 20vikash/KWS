package store

import (
	"context"
	"errors"
	"kws/kws/consts/status"
	"kws/kws/models"
	"kws/kws/models/web"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InstanceStore struct {
	db *pgxpool.Pool
}

// Database level function: Insert an instance record related to the user
func (i *InstanceStore) CreateInstance(ctx context.Context, uid int, userName, insUser, insPassword string) error {
	instance := models.CreateInstanceType(uid, userName)

	sql := `
		INSERT INTO instance (user_id, volume_name, container_name, instance_type, is_running, ins_user, ins_password) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := i.db.Exec(ctx, sql,
		instance.Uid,
		instance.VolumeName,
		instance.ContainerName,
		instance.InstanceType,
		true,
		insUser,
		insPassword,
	)
	if err != nil {
		log.Println("Cannot insert row into instance table")
		return err
	}

	return nil
}

// Database level function: Start an instance record related to the user
func (i *InstanceStore) StartInstance(ctx context.Context, uid int) error {
	var isRunning bool

	sql := `
		SELECT is_running FROM instance WHERE user_id=$1
	`

	err := i.db.QueryRow(ctx, sql, uid).Scan(&isRunning)
	if err != nil {
		log.Println("Cannot query the instance table(db)")
		return err
	}

	if isRunning {
		return errors.New(status.CONTAINER_ALREADY_RUNNING)
	}

	sql = `
		UPDATE instance SET is_running=TRUE WHERE user_id=$1
	`

	res, err := i.db.Exec(ctx, sql, uid)
	if err != nil {
		log.Println("Cannot start the instance(db)")
		return err
	}

	rows := res.RowsAffected()

	if rows == 0 {
		log.Printf("No user found with user_id=%d\n", uid)
		return errors.New(status.CONTAINER_START_FAILED)
	}

	return nil
}

// Database level function: Remove an instance record related to the user
func (i *InstanceStore) RemoveInstance(ctx context.Context, uid int) error {
	sql := `
		DELETE FROM instance WHERE user_id = $1
	`

	res, err := i.db.Exec(ctx, sql, uid)
	if err != nil {
		log.Println("Cannot remove the instance(db)")
		return err
	}

	rows := res.RowsAffected()

	if rows == 0 {
		log.Printf("No user found with user_id=%d\n", uid)
		return errors.New(status.CONTAINER_STOP_FAILED)
	}

	return nil
}

// Database level function: Stop an instance record related to the user
func (i *InstanceStore) StopInstance(ctx context.Context, uid int) error {
	sql := `
		UPDATE instance SET is_running=FALSE WHERE user_id=$1
	`

	res, err := i.db.Exec(ctx, sql, uid)
	if err != nil {
		log.Println("Cannot stop the instance(db)")
		return err
	}

	rows := res.RowsAffected()

	if rows == 0 {
		log.Printf("No user found with user_id=%d\n", uid)
		return errors.New(status.CONTAINER_DELETE_FAILED)
	}

	return nil
}

func (i *InstanceStore) Exists(ctx context.Context, uid int) (bool, error) {
	var c int

	sql := `
		SELECT COUNT(id) FROM instance WHERE user_id  = $1
	`

	err := i.db.QueryRow(ctx, sql, uid).Scan(&c)
	if err != nil {
		log.Println("Failed to check instance existance")
		return false, err
	}

	if c != 1 {
		return false, nil
	}

	return true, nil
}

func (i *InstanceStore) GetData(ctx context.Context, uid int) (*web.InsData, error) {
	var insData = new(web.InsData)
	var isRunning bool

	sql := `
		SELECT ins_user, ins_password, container_name, is_running
		FROM instances
		WHERE user_id = $1
	`

	err := i.db.QueryRow(ctx, sql, uid).Scan(
		&insData.Instance.Username,
		&insData.Instance.Password,
		&insData.ContainerName,
		&isRunning,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			insData.InstanceStatus = "inactive"
			insData.Active = "no"
			return insData, nil
		}

		log.Println("Cannot find instance data:", err)
		return nil, err
	}

	if isRunning {
		insData.InstanceStatus = "active"
		insData.Active = "exists"
	} else {
		insData.InstanceStatus = "stopped"
		insData.Active = "no"
	}

	return insData, nil
}
