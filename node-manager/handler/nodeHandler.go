package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"time"

	"github.com/dipghoshraj/media-service/node-manager/model"
	"github.com/go-redis/redis/v8"
)

type NodeStatus string

const (
	NodeStatusActive  NodeStatus = "active"
	NodeStatusFailed  NodeStatus = "failed"
	NodeStatusPending NodeStatus = "pending"
)

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

func (ndb *DBHandler) RegisterNode(node *model.Node) error {
	// Implement the logic to register a node
	// Return an error if the node registration fails

	tx := ndb.DbManager.DB.Begin()

	if err := tx.Create(node).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to store node in database: %v", err)
	}

	nodeJSON, err := json.Marshal(node)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to marshal node: %v", err)
	}

	err = ndb.DbManager.RedisClient.Set(context.Background(),
		fmt.Sprintf("node:%s", node.ID),
		nodeJSON,
		24*time.Hour).Err()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to store node in Redis: %v", err)
	}

	// SpinUpContainer
	port, err := ndb.SpinUpContainer(node.ID, node.VolumeName, node.Capacity)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to spin up the container: %v", err)
	}

	node.Port = port
	if err := tx.Save(node).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update node port: %v", err)
	}
	return tx.Commit().Error
}

func (ndb *DBHandler) FindNode() (model.Node, error) {
	// Implement the logic to find a node
	// Return the node

	var node model.Node
	nodevalue := ndb.DbManager.DB.Where("status = ? AND capacity - used_space >= ?", "active", 5).First(&node)
	if nodevalue.Error != nil {
		return node, fmt.Errorf("failed to find node: %v", nodevalue.Error)
	}
	return node, nil
}

func (ndb *DBHandler) SpinUpContainer(nodeID string, volumName string, capacity int64) (string, error) {

	if !volumeExists(volumName) {
		if err := createVolume(volumName); err != nil {
			return "", fmt.Errorf("error creating volume: %v", err)
		}
	}

	port, err := AllocatePort(ndb.DbManager.RedisClient)
	if err != nil {
		return "", fmt.Errorf("failed to port allocation: %v", err)
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
		return "", fmt.Errorf("failed to run docker: %v", err)
	}
	fmt.Println("Container started:", string(output))
	return port, nil
}

func (ndb *DBHandler) GetClusterStats() (map[string]interface{}, error) {
	var totalCapacity, totalUsedSpace int64
	var activeNodes, totalNodes int64

	// Get aggregate statistics
	err := ndb.DbManager.DB.Model(&model.Node{}).
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

func (ndb *DBHandler) GetAllNodes() ([]model.Node, error) {
	var nodes []model.Node
	if err := ndb.DbManager.DB.Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}
