package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"qixue-api/internal/platform/config"
)

func Open(cfg config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis: %w", err)
	}
	return rdb, nil
}
