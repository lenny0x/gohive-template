package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func Init(cfg Config) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	return err
}

func GetClient() *redis.Client {
	return Client
}
