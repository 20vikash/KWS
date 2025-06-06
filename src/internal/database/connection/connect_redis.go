package database

import "github.com/redis/go-redis/v9"

type RedisDB struct {
	Addr     string
	Password string
	DB       int
}

func (r *RedisDB) Connect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})

	return rdb
}
