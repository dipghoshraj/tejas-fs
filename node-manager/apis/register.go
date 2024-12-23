package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dipghoshraj/media-service/node-manager/model"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type NodeRegistrationRequest struct {
	IP       string `json:"ip"`
	Capacity int64  `json:"capacity"`
}

type HeartbeatRequest struct {
	UsedSpace int64 `json:"usedSpace"`
}

type NMHandler struct {
	NodeManager *model.NodeManager
}

func SetupRoutes(router *mux.Router, nm *NMHandler) {
	// Node management endpoints
	router.HandleFunc("/api/nodes", nm.RegisterNodeHandler).Methods("POST")
	router.HandleFunc("/api/nodes/stats", nm.GetClusterStatsHandler).Methods("GET")
	router.HandleFunc("/api/all/nodes", nm.GetAllNodesHandler).Methods("GET")

	router.Use(loggingMiddleware)
	router.Use(recoveryMiddleware)

}

func NewNMHandler(NodeManager *model.NodeManager) *NMHandler {
	return &NMHandler{NodeManager: NodeManager}
}

func (nm *NMHandler) RegisterNodeHandler(w http.ResponseWriter, r *http.Request) {
	var request NodeRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to decode request: %v", err),
		})
		return
	}

	ip, err := nm.NodeManager.AllocateIP()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to allocate IP: %v", err),
		})
		return
	}

	node := &model.Node{
		ID:            fmt.Sprintf("node-%d", time.Now().Unix()),
		IP:            ip,
		Status:        model.NodeStatusActive,
		Capacity:      request.Capacity,
		UsedSpace:     0,
		LastHeartbeat: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := nm.NodeManager.CreateNode(node); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to register node: %v", err),
		})
		return
	}

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    node,
	})
}

func (nm *NMHandler) GetAllNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := nm.NodeManager.GetAllNodes()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch nodes")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    nodes,
	})
}

func (nm *NMHandler) GetClusterStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := nm.NodeManager.GetClusterStats()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch nodes")
		return
	}
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		Success: false,
		Error:   message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				respondWithError(w, http.StatusInternalServerError, "Internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
