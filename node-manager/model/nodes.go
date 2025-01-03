package model

import (
	"time"
)

type NodeStatus string

const (
	NodeStatusActive   NodeStatus = "active"
	NodeStatusFailed   NodeStatus = "failed"
	NodeStatusPending  NodeStatus = "pending"
	NodeStatusInactive NodeStatus = "inactive"
)

type Node struct {
	ID            string     `json:"id" gorm:"primaryKey"`
	Status        NodeStatus `json:"status"`
	Capacity      int64      `json:"capacity"`  // in bytes
	UsedSpace     int64      `json:"usedSpace"` // in bytes
	LastHeartbeat time.Time  `json:"lastHeartbeat"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	VolumeName    string     `json:"volumeName"`
	Port          string     `json:"port"`
}

// type NodeManager struct {
// 	DB          *gorm.DB
// 	redisClient *redis.Client
// 	Lock        sync.Mutex
// 	ctx         context.Context
// }
