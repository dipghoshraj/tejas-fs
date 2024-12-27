package model

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type DbManager struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	Lock        sync.Mutex
	ctx         context.Context
}

func NewNodeManager(db *gorm.DB, redisClient *redis.Client, ipPool string) *DbManager {
	ctx := context.Background()

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil
	}

	return &DbManager{
		DB:          db,
		RedisClient: redisClient,
		ctx:         context.Background(),
	}
}
