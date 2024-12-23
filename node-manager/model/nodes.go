package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
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
	DB          *gorm.DB
	redisClient *redis.Client
	Nodes       map[string]*Node
	Lock        sync.Mutex
	ipPool      []string
	usedIPs     map[string]bool
	ctx         context.Context
}

func NewNodeManager(db *gorm.DB, redisClient *redis.Client) *NodeManager {
	return &NodeManager{
		DB:          db,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// ClusterManager manages the IP address pool
func (cm *NodeManager) AllocateIP() (string, error) {
	for _, ip := range cm.ipPool {
		if !cm.usedIPs[ip] {
			cm.usedIPs[ip] = true
			return ip, nil
		}
	}
	return "", fmt.Errorf("no available IPs")
}

// ReleaseIP releases an IP address back to the pool
func (cm *NodeManager) ReleaseIP(ip string) {
	cm.Lock.Lock()
	defer cm.Lock.Unlock()
	delete(cm.usedIPs, ip)
}

func (nm *NodeManager) CreateNode(node *Node) error {

	if err := nm.DB.Create(node).Error; err != nil {
		nm.DB.Rollback()
		return fmt.Errorf("failed to store node in database: %v", err)
	}

	nodeJSON, err := json.Marshal(node)
	if err != nil {
		nm.DB.Rollback()
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	err = nm.redisClient.Set(context.Background(),
		fmt.Sprintf("node:%s", node.ID),
		nodeJSON,
		24*time.Hour).Err()
	if err != nil {
		nm.DB.Rollback()
		return fmt.Errorf("failed to store node in Redis: %v", err)
	}

	return nm.DB.Commit().Error

}

// UpdateNodeStatus updates the status of a node
func (nm *NodeManager) UpdateNodeStatus(nodeID string, status NodeStatus) error {
	node, exists := nm.Nodes[nodeID]
	if !exists {
		return fmt.Errorf("node not found: %s", nodeID)
	}

	node.Status = status
	node.UpdatedAt = time.Now()

	// Update in PostgreSQL
	if err := nm.DB.Save(node).Error; err != nil {
		nm.DB.Rollback()
		return fmt.Errorf("failed to update node in database: %v", err)
	}

	// Update in Redis
	nodeJSON, err := json.Marshal(node)
	if err != nil {
		nm.DB.Rollback()
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	err = nm.redisClient.Set(context.Background(),
		fmt.Sprintf("node:%s", node.ID),
		nodeJSON,
		24*time.Hour).Err()
	if err != nil {
		nm.DB.Rollback()
		return fmt.Errorf("failed to update node in Redis: %v", err)
	}

	return nil
}

// HandleHeartbeat processes a heartbeat from a node
func (nm *NodeManager) HandleHeartbeat(nodeID string) error {
	node, exists := nm.Nodes[nodeID]
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
			for _, node := range nm.Nodes {
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
func (nm *NodeManager) GetClusterStats() (map[string]interface{}, error) {
	var totalCapacity, totalUsedSpace int64
	var activeNodes, totalNodes int64

	// Get aggregate statistics
	err := nm.DB.Model(&Node{}).
		Select("COUNT(*) as total_nodes, "+
			"SUM(capacity) as total_capacity, "+
			"SUM(used_space) as total_used_space, "+
			"COUNT(CASE WHEN status = ? THEN 1 END) as active_nodes",
			NodeStatusActive).
		Row().
		Scan(&totalNodes, &totalCapacity, &totalUsedSpace, &activeNodes)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"totalNodes":     totalNodes,
		"activeNodes":    activeNodes,
		"totalCapacity":  totalCapacity,
		"totalUsedSpace": totalUsedSpace,
		"usagePercent":   float64(totalUsedSpace) / float64(totalCapacity) * 100,
	}, nil

}

func (ops *NodeManager) GetAllNodes() ([]Node, error) {
	var nodes []Node
	if err := ops.DB.Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}
