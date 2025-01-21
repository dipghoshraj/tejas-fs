package handler

import (
	cosmicmodel "planet-of-node/cosmic-model"
)

// Init initializes the cosmic model
type hManager struct {
	dbm *cosmicmodel.DBM
}

func HandlerManager(dbm *cosmicmodel.DBM) *hManager {
	return &hManager{dbm: dbm}
}
