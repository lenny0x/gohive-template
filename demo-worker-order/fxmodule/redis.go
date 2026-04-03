package fxmodule

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/fx"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-order/config"
)

func NewRedis(lc fx.Lifecycle) *goredis.Client {
	cfg := config.Cfg.Redis
	client := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if _, err := client.Ping(ctx).Result(); err != nil {
				return fmt.Errorf("redis ping failed: %w", err)
			}
			logger.Info("Redis connected")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Redis connection closed")
			return client.Close()
		},
	})

	return client
}
