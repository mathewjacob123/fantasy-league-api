package db

import (
	"fantasy-league-api/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) *redis.Client {
	if cfg.RedisHost == "" {
		log.Println("Redis not configured — caching disabled")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
	})

	return client
}