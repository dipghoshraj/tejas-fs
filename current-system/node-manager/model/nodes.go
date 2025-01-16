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
	Capacity      int64      `json:"capacity"`  // in MB
	UsedSpace     int64      `json:"usedSpace"` // in MB
	LastHeartbeat time.Time  `json:"lastHeartbeat"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	VolumeName    string     `json:"volumeName"`
	Port          string     `json:"port"`
}
