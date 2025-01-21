package nodeunit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func AllocatePort(rdb *redis.Client) (string, error) {
	deadline := time.Now().Add(20 * time.Second) // 20 seconds timeout
	ctx := context.Background()

	for {
		port, err := rdb.SPop(ctx, "available_ports").Result()
		if err == nil {
			return port, nil
		}

		if time.Now().After(deadline) {
			return "", fmt.Errorf("timeout reached, no available ports")
		}

		if err != redis.Nil {
			return "", fmt.Errorf("error fetching port: %w", err)
		}

		fmt.Println("No available ports, retrying...")
		time.Sleep(100 * time.Millisecond) // Small delay before retrying
	}
}
