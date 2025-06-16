package store

import (
	"context"
	"errors"
	"kws/kws/consts/config"
	"kws/kws/consts/status"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	ds *redis.Client
}

func (r *RedisStore) SetEmailToken(ctx context.Context, email string, token string) error {
	err := r.ds.Set(ctx, token, "email:"+email, 24*time.Hour).Err()
	if err != nil {
		log.Println(err.Error())
		log.Println("Failed to set email token")
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

func (r *RedisStore) PushFreeIp(ctx context.Context, ip int) error {
	err := r.ds.LPush(ctx, config.STACK_KEY, ip).Err()
	if err != nil {
		log.Println("Cannot push free IP to the stack")
		return err
	}

	return nil
}

func (r *RedisStore) PopFreeIp(ctx context.Context) (int, error) {
	val, err := r.ds.LPop(ctx, config.STACK_KEY).Result()
	if err != nil {
		if err == redis.Nil {
			log.Println("Nothing to pop")
			return -1, errors.New(status.EMPTY_IP_STACK)
		}

		return -1, err
	}

	intVal, convErr := strconv.Atoi(val)
	if convErr != nil {
		log.Println("Conversion to int failed")
		return -1, convErr
	}

	return intVal, nil
}
