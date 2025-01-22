package handler

import (
	cosmicmodel "planet-of-node/cosmic-model"
)

func (hm *hManager) GetClusterStats() (map[string]interface{}, error) {
	var totalCapacity, totalUsedSpace int64
	var activeNodes, totalNodes int64

	// Get aggregate statistics
	err := hm.dbm.DB.Model(&cosmicmodel.Node{}).
		Select("COUNT(*) as total_nodes, "+
			"SUM(capacity) as total_capacity, "+
			"SUM(used_space) as total_used_space, "+
			"COUNT(CASE WHEN status = ? THEN 1 END) as active_nodes",
			cosmicmodel.NodeStatusActive).
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
