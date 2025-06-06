package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Auth interface {
		CreateUser(ctx context.Context, first_name, last_name, email, password, user_name string) error
	}

	InMemory interface {
		SetEmailToken(ctx context.Context, email string, token string) error
		GetEmailFromToken(ctx context.Context, token string) string
		DeleteEmailToken(ctx context.Context, token string) error
	}
}

func NewStore(pg *pgxpool.Pool, redis *redis.Client) *Storage {
	return &Storage{
		Auth: &AuthStore{
			db: pg,
		},
		InMemory: &RedisStore{
			ds: redis,
		},
	}
}
