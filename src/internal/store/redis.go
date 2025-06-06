package store

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	ds *redis.Client
}

func (r *RedisStore) SetEmailToken(ctx context.Context, email string, token string) error {
	err := r.ds.Set(ctx, token, "email:"+email, 24*time.Hour).Err()
	if err != nil {
		log.Fatal("Failed to set email token")
		return err
	}

	return nil
}

func (r *RedisStore) DeleteEmailToken(ctx context.Context, token string) error {
	res, _ := r.ds.Del(ctx, token).Result()
	if res == 0 {
		return errors.New("key expired")
	}

	return nil
}

func (r *RedisStore) GetEmailFromToken(ctx context.Context, token string) string {
	val := r.ds.Get(ctx, token).String()

	return val
}
