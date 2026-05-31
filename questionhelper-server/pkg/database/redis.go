package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"questionhelper-server/pkg/config"
)

var RDB *redis.Client

func InitRedis(cfg config.RedisConfig) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("连接 Redis 失败: %v", err))
	}
}
