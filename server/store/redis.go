package store

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type Redis struct {
	Client *redis.Client
}

func NewRedis() *Redis {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	username := os.Getenv("REDIS_UNAME")
	password := os.Getenv("REDIS_PASS")
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       0,
	})
	return &Redis{Client: rdb}
}

func (r *Redis) SetJTI(ctx context.Context, key, userID string, exp time.Time) error {
	return r.Client.Set(ctx, key, userID, time.Until(exp)).Err()
}

func (r *Redis) DelJTI(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) GetUserByJTI(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
