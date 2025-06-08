package store

import (
	"context"
	"kws/kws/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Auth interface {
		CreateUser(ctx context.Context, first_name, last_name, email, password, user_name string) error
		GenerateToken(ctx context.Context, email string) string
		VerifyUser(ctx context.Context, email string) error
		LoginUser(ctx context.Context, userName, password string) (*models.User, error)
	}

	InMemory interface {
		SetEmailToken(ctx context.Context, email string, token string) error
		GetEmailFromToken(ctx context.Context, token string) string
		DeleteEmailToken(ctx context.Context, token string) error
	}

	Instance interface {
		CreateInstance(ctx context.Context, uid int, userName string) error
		RemoveInstance(ctx context.Context, uid int) error
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
		Instance: &InstanceStore{
			db: pg,
		},
	}
}
