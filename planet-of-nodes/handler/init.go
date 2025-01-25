package handler

import (
	cosmicmodel "planet-of-node/cosmic-model"
)

// Init initializes the cosmic model
type HManager struct {
	dbm *cosmicmodel.DBM
}

func HandlerManager(dbm *cosmicmodel.DBM) *HManager {
	return &HManager{dbm: dbm}
}
