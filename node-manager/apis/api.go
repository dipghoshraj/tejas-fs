package apis

import (
	"github.com/dipghoshraj/media-service/node-manager/model"
	"github.com/gorilla/mux"
)

type NMHandler struct {
	DbManager *model.DbManager
}

func SetupRoutes(router *mux.Router, nm *NMHandler) {
	// Node management endpoints
	router.HandleFunc("/api/nodes", nm.RegisterNodeHandler).Methods("POST")
	router.HandleFunc("/api/nodes/stats", nm.GetClusterStatsHandler).Methods("GET")
	router.HandleFunc("/api/all/nodes", nm.GetAllNodesHandler).Methods("GET")

	router.Use(loggingMiddleware)
	router.Use(recoveryMiddleware)
}
