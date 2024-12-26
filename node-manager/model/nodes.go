package model

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type NodeStatus string

const (
	NodeStatusActive  NodeStatus = "active"
	NodeStatusFailed  NodeStatus = "failed"
	NodeStatusPending NodeStatus = "pending"
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
}

// type NodeManager struct {
// 	DB          *gorm.DB
// 	redisClient *redis.Client
// 	Lock        sync.Mutex
// 	ctx         context.Context
// }

func (nm *DbManager) CreateNode(node *Node) error {
	tx := nm.DB.Begin()

	if err := tx.Create(node).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to store node in database: %v", err)
	}

	nodeJSON, err := json.Marshal(node)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	err = nm.redisClient.Set(context.Background(),
		fmt.Sprintf("node:%s", node.ID),
		nodeJSON,
		24*time.Hour).Err()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to store node in Redis: %v", err)
	}

	// SpinUpContainer
	if err := nm.SpinUpContainer(node.ID, node.VolumeName, node.Capacity); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to spin up the container: %v", err)
	}

	return tx.Commit().Error

}

// UpdateNodeStatus updates the status of a node
// func (nm *NodeManager) UpdateNodeStatus(nodeID string, status NodeStatus) error {
// 	node, exists := nm.Nodes[nodeID]
// 	if !exists {
// 		return fmt.Errorf("node not found: %s", nodeID)
// 	}

// 	node.Status = status
// 	node.UpdatedAt = time.Now()

// 	// Update in PostgreSQL
// 	if err := nm.DB.Save(node).Error; err != nil {
// 		nm.DB.Rollback()
// 		return fmt.Errorf("failed to update node in database: %v", err)
// 	}

// 	// Update in Redis
// 	nodeJSON, err := json.Marshal(node)
// 	if err != nil {
// 		nm.DB.Rollback()
// 		return fmt.Errorf("failed to marshal node: %v", err)
// 	}

// 	err = nm.redisClient.Set(context.Background(),
// 		fmt.Sprintf("node:%s", node.ID),
// 		nodeJSON,
// 		24*time.Hour).Err()
// 	if err != nil {
// 		nm.DB.Rollback()
// 		return fmt.Errorf("failed to update node in Redis: %v", err)
// 	}

// 	return nil
// }

// HandleHeartbeat processes a heartbeat from a node
// func (nm *NodeManager) HandleHeartbeat(nodeID string) error {
// 	node, exists := nm.Nodes[nodeID]
// 	if !exists {
// 		return fmt.Errorf("node not found: %s", nodeID)
// 	}

// 	node.LastHeartbeat = time.Now()
// 	return nm.UpdateNodeStatus(nodeID, NodeStatusActive)

// }

// MonitorNodes periodically checks node health
// func (nm *NodeManager) MonitorNodes(timeout time.Duration) {
// 	ticker := time.NewTicker(30 * time.Second)
// 	go func() {
// 		for range ticker.C {
// 			now := time.Now()
// 			for _, node := range nm.Nodes {
// 				if now.Sub(node.LastHeartbeat) > timeout {
// 					log.Printf("Node %s missed heartbeat, marking as failed", node.ID)
// 					err := nm.UpdateNodeStatus(node.ID, NodeStatusFailed)
// 					if err != nil {
// 						log.Printf("Failed to update node status: %v", err)
// 					}
// 				}
// 			}
// 		}
// 	}()
// }

// GetNodeStats returns current statistics about the node cluster
func (nm *DbManager) GetClusterStats() (map[string]interface{}, error) {
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

func (ops *DbManager) GetAllNodes() ([]Node, error) {
	var nodes []Node
	if err := ops.DB.Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (nm *DbManager) SpinUpContainer(nodeID string, volumName string, capacity int64) error {

	if !volumeExists(volumName) {
		if err := createVolume(volumName); err != nil {
			return fmt.Errorf("error creating volume: %v", err)
		}
	}

	port, err := nm.AllocatePort(nm.redisClient)
	if err != nil {
		return fmt.Errorf("failed to port allocation: %v", err)
	}

	cmd := exec.Command("docker", "run", "-d",
		"--name", nodeID,
		"--env", fmt.Sprintf("NODE_ID=%s", nodeID),
		"--env", fmt.Sprintf("STORAGE_CAPACITY=%s", strconv.FormatInt(capacity, 10)),
		"-v", fmt.Sprintf("%s:/data", volumName),
		"-p", fmt.Sprintf("%s:%s", port, "8080"),
		"node-image")

	fmt.Println("Executing command:", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Output:", string(output))
		return fmt.Errorf("failed to run docker: %v", err)
	}
	fmt.Println("Container started:", string(output))
	return nil
}

func volumeExists(volumeName string) bool {
	cmd := exec.Command("docker", "volume", "inspect", volumeName)
	err := cmd.Run()
	return err == nil // If there is no error the volume exist
}

func createVolume(volumeName string) error {
	cmd := exec.Command("docker", "volume", "create", volumeName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating volume: %w, output: %s", err, string(output))
	}
	fmt.Println("volume created", string(output))
	return nil
}

func (cm *DbManager) AllocatePort(rdb *redis.Client) (string, error) {
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

// func releasePort(rdb *redis.Client, port int) {
// 	ctx := context.Background()

// 	rdb.SAdd(ctx, "available_ports", port)
// 	fmt.Printf("Port %d released\n", port)
// }

// func incrementIP(ip net.IP) {
// 	for j := len(ip) - 1; j >= 0; j-- {
// 		ip[j]++
// 		if ip[j] > 0 {
// 			break
// 		}
// 	}
// }

// ClusterManager manages the IP address pool
// func (cm *NodeManager) AllocateIP() (string, error) {
// 	for _, ip := range cm.ipPool {
// 		key := fmt.Sprintf("ip:%s", ip)
// 		// Use SETNX (SET if Not eXists) for atomic allocation
// 		set, err := cm.redisClient.SetNX(cm.ctx, key, "used", 0).Result()
// 		if err != nil {
// 			log.Fatalf("failed to allocate IP in Redis: %v", err)
// 			continue
// 		}
// 		if set { // If the key was set, it means the IP was free
// 			return ip, nil
// 		}
// 	}

// 	return "", fmt.Errorf("no available IPs")
// }

// ReleaseIP releases an IP address back to the pool
// func (cm *NodeManager) ReleaseIP(ip string) error {
// 	key := fmt.Sprintf("ip:%s", ip)
// 	err := cm.redisClient.Del(cm.ctx, key).Err()
// 	return err
// }
