package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type NodeStatus string

const (
	NodeStatusActive  NodeStatus = "active"
	NodeStatusFailed  NodeStatus = "failed"
	NodeStatusPending NodeStatus = "pending"
)

type Node struct {
	ID            string     `json:"id" gorm:"primaryKey"`
	IP            string     `json:"ip"`
	Status        NodeStatus `json:"status"`
	Capacity      int64      `json:"capacity"`  // in bytes
	UsedSpace     int64      `json:"usedSpace"` // in bytes
	LastHeartbeat time.Time  `json:"lastHeartbeat"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

type NodeManager struct {
	db          *gorm.DB
	redisClient *redis.Client
	nodes       map[string]*Node
}

func (nm *NodeManager) RegisterNode(node *Node) error {

	if err := nm.db.Create(node).Error; err != nil {
		return fmt.Errorf("failed to store node in database: %v", err)
	}

	nodeJSON, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	err = nm.redisClient.Set(context.Background(),
		fmt.Sprintf("node:%s", node.ID),
		nodeJSON,
		24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store node in Redis: %v", err)
	}

	// Store in memory
	nm.nodes[node.ID] = node
	return nil
}

// UpdateNodeStatus updates the status of a node
func (nm *NodeManager) UpdateNodeStatus(nodeID string, status NodeStatus) error {
	node, exists := nm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.Status = status
	node.UpdatedAt = time.Now()

	// Update in PostgreSQL
	if err := nm.db.Save(node).Error; err != nil {
		return fmt.Errorf("failed to update node in database: %v", err)
	}

	// Update in Redis
	nodeJSON, err := json.Marshal(node)
	if err != nil {
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	err = nm.redisClient.Set(context.Background(),
		fmt.Sprintf("node:%s", node.ID),
		nodeJSON,
		24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to update node in Redis: %v", err)
	}

	return nil
}

// HandleHeartbeat processes a heartbeat from a node
func (nm *NodeManager) HandleHeartbeat(nodeID string) error {
	node, exists := nm.nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.LastHeartbeat = time.Now()
	return nm.UpdateNodeStatus(nodeID, NodeStatusActive)

}

// MonitorNodes periodically checks node health
func (nm *NodeManager) MonitorNodes(timeout time.Duration) {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			now := time.Now()
			for _, node := range nm.nodes {
				if now.Sub(node.LastHeartbeat) > timeout {
					log.Printf("Node %s missed heartbeat, marking as failed", node.ID)
					err := nm.UpdateNodeStatus(node.ID, NodeStatusFailed)
					if err != nil {
						log.Printf("Failed to update node status: %v", err)
					}
				}
			}
		}
	}()
}

// GetNodeStats returns current statistics about the node cluster
func (nm *NodeManager) GetNodeStats() map[string]interface{} {
	totalNodes := len(nm.nodes)
	activeNodes := 0
	totalCapacity := int64(0)
	totalUsedSpace := int64(0)

	for _, node := range nm.nodes {
		if node.Status == NodeStatusActive {
			activeNodes++
			totalCapacity += node.Capacity
			totalUsedSpace += node.UsedSpace
		}
	}

	return map[string]interface{}{
		"totalNodes":     totalNodes,
		"activeNodes":    activeNodes,
		"totalCapacity":  totalCapacity,
		"totalUsedSpace": totalUsedSpace,
	}
}

// type Config struct {
// 	PostgresURL      string
// 	RedisURL         string
// 	HeartbeatTimeout time.Duration
// }

// func DbConnect(config Config) (*NodeManager, error) {
// 	// Initialize PostgreSQL connection
// 	db, err := gorm.Open(postgres.Open(config.PostgresURL), &gorm.Config{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect to database: %v", err)
// 	}

// 	// Auto-migrate the schema
// 	err = db.AutoMigrate(&Node{})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to migrate database: %v", err)
// 	}

// 	// Initialize Redis client
// 	redisClient := redis.NewClient(&redis.Options{
// 		Addr: config.RedisURL,
// 	})

// 	return &NodeManager{
// 		db:          db,
// 		redisClient: redisClient,
// 		nodes:       make(map[string]*Node),
// 	}, nil
// }
