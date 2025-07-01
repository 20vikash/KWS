package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kws/kws/consts/status"
	"kws/kws/models/web"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Ds *redis.Client
}

// Email
func (r *RedisStore) SetEmailToken(ctx context.Context, email string, token string) error {
	err := r.Ds.Set(ctx, token, "email:"+email, 24*time.Hour).Err()
	if err != nil {
		log.Println(err.Error())
		log.Println("Failed to set email token")
		return err
	}

	return nil
}

func (r *RedisStore) DeleteEmailToken(ctx context.Context, token string) error {
	res, _ := r.Ds.Del(ctx, token).Result()
	if res == 0 {
		return errors.New("key expired")
	}

	return nil
}

func (r *RedisStore) GetEmailFromToken(ctx context.Context, token string) string {
	val := r.Ds.Get(ctx, token).String()

	return val
}

// Wireguard
func (r *RedisStore) PushFreeIp(ctx context.Context, ip int, key string) error {
	err := r.Ds.LPush(ctx, key, ip).Err()
	if err != nil {
		log.Println("Cannot push free IP to the stack")
		return err
	}

	return nil
}

func (r *RedisStore) PopFreeIp(ctx context.Context, key string) (int, error) {
	val, err := r.Ds.LPop(ctx, key).Result()
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

// Instance deploy, kill and stop

func (r *RedisStore) PutDeployResult(ctx context.Context, userName, jobID, password, ip string, success bool) error {
	instance := web.Instance{
		Success:  success,
		Username: userName,
		Password: password,
		IP:       ip,
	}

	data, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	return r.Ds.Set(ctx, "deploy:result:"+jobID, data, 0).Err()
}

func (r *RedisStore) GetDeployResult(ctx context.Context, jobID string) (bool, *web.Instance, error) {
	val, err := r.Ds.Get(ctx, "deploy:result:"+jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil, nil
		}
		return false, nil, err
	}

	var instance web.Instance
	if err := json.Unmarshal([]byte(val), &instance); err != nil {
		return false, nil, err
	}

	return true, &instance, nil
}

func (r *RedisStore) PutStopResult(ctx context.Context, result bool, jobID string) error {
	return r.Ds.Set(ctx, "stop:result:"+jobID, fmt.Sprintf("%v", result), 0).Err()
}

func (r *RedisStore) GetStopResult(ctx context.Context, jobID string) (bool, bool, error) {
	val, err := r.Ds.Get(ctx, "stop:result:"+jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, false, nil
		}
		return false, false, err
	}

	return true, val == "true", nil
}

func (r *RedisStore) PutKillResult(ctx context.Context, result bool, jobID string) error {
	return r.Ds.Set(ctx, "kill:result:"+jobID, fmt.Sprintf("%v", result), 0).Err()
}

func (r *RedisStore) GetKillResult(ctx context.Context, jobID string) (bool, bool, error) {
	val, err := r.Ds.Get(ctx, "kill:result:"+jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return false, false, nil
		}
		return false, false, err
	}

	return true, val == "true", nil
}
