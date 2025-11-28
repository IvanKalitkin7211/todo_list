package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
	"todo-list/config"
)

func ProvideRedisClient(cfg *config.RedisConfig) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	}
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Printf("[ERROR] redis ping: %v", err)
		return nil, err
	}
	log.Printf("[INFO] redis connected")
	return client, nil
}
