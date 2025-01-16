package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dipghoshraj/media-service/file-service-nodes/infrastructure"
	"github.com/go-redis/redis/v8"
)

func InitializeRedis(config *infrastructure.Infrastructure) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return redisClient, nil
}
