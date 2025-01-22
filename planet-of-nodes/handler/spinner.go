package handler

import (
	"fmt"
	cosmicmodel "planet-of-node/cosmic-model"
	nodeunit "planet-of-node/node-unit"
)

type NodeStatus string

const (
	NodeStatusActive  NodeStatus = "active"
	NodeStatusFailed  NodeStatus = "failed"
	NodeStatusPending NodeStatus = "pending"
)

func (hm *hManager) SpinUpContainer(node *cosmicmodel.Node) error {
	// start the transcation operation
	tx := hm.dbm.DB.Begin()
	// crate the node metadata in theDB
	if err := tx.Create(node).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update node port: %v", err)
	}

	port, err := nodeunit.AllocatePort(hm.dbm.RedisClient)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update node port: %v", err)
	}

	node.Port = port

	if err := nodeunit.SpinUpContainer(node.ID, node.VolumeName, node.Capacity, node.Port); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update node port: %v", err)
	}

	if err := tx.Save(node).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update node port: %v", err)
	}

	return nil
}
