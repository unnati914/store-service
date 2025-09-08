package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func ConnectRedis(addr string) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Redis{Client: rdb}, nil
}

func (r *Redis) Close() error { return r.Client.Close() }
