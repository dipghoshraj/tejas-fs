package handler

import (
	"fmt"
	cosmicmodel "planet-of-node/cosmic-model"
)

func (hm *HManager) CreateCusterMetadata(clsuter *cosmicmodel.ClusterConsfig) error {

	tx := hm.dbm.DB.Begin()
	if err := tx.Create(clsuter).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create custer metadata: %w", err)
	}
	return nil
}
