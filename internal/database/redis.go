package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/sayeed1999/ride-sharing-golang-api/config"
)

func InitRedis(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func InitRedisWithErrorHandling(cfg *config.Config) *redis.Client {
	client, err := InitRedis(cfg)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	return client
}

func CloseRedis(client *redis.Client) error {
	return client.Close()
}
