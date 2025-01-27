package handler

import (
	"context"
	"fmt"
	cosmicmodel "planet-of-node/cosmic-model"
	"sync"
	"time"

	"github.com/google/uuid"
)

func (hm *HManager) CreateCusterMetadata(cluster *cosmicmodel.Cluster) error {

	tx := hm.dbm.DB.Begin()
	if err := tx.Create(cluster).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create custer metadata: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	go hm.NodeCreateParallel(ctx, cluster) // runn the node creation in background
	defer cancel()
	return tx.Commit().Error
}

func (hm *HManager) NodeCreateParallel(ctx context.Context, cluster *cosmicmodel.Cluster) error {
	var wg sync.WaitGroup

	done := make(chan struct{})
	defer close(done)

	for iter := 0; iter < int(cluster.Nodes); iter++ {
		wg.Add(1)

		go func(iter int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				node := &cosmicmodel.Node{
					ID:            uuid.New().String(),
					VolumeName:    uuid.New().String(),
					Status:        cosmicmodel.NodeStatusInactive,
					Capacity:      cluster.NodeCapacity,
					UsedSpace:     0,
					LastHeartbeat: time.Now(),
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
					Cluster_id:    cluster.ID,
					Host:          "localhost",
				}

				if err := hm.SpinUpContainer(node); err != nil {
					fmt.Printf("Error creating node %d: %v\n", iter, err)
					return
				}
			}
		}(iter)
	}

	// go func() {
	// 	wg.Wait()
	// 	close(done)
	// }()

	<-ctx.Done()
	return fmt.Errorf("cluster spinup errors %v", ctx.Err())
}
