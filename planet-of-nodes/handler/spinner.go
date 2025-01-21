package handler

import (
	cosmicmodel "planet-of-node/cosmic-model"
)

type NodeStatus string

const (
	NodeStatusActive  NodeStatus = "active"
	NodeStatusFailed  NodeStatus = "failed"
	NodeStatusPending NodeStatus = "pending"
)

func (hm *hManager) SpinUpContainer(node *cosmicmodel.Node) error {
	// start the transcation operation

	return nil
}
