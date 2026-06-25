package db

import (
	"fantasy-league-api/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	return client
}