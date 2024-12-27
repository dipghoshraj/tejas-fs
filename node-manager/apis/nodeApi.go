package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dipghoshraj/media-service/node-manager/handler"
	"github.com/dipghoshraj/media-service/node-manager/model"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type NodeRegistrationRequest struct {
	Capacity int64 `json:"capacity"`
}

type HeartbeatRequest struct {
	UsedSpace int64 `json:"usedSpace"`
}

func NewNMHandler(NodeManager *handler.DBHandler) *NMHandler {
	return &NMHandler{DbManager: NodeManager}
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

	volumnName := uuid.New().String()

	node := &model.Node{
		ID:            fmt.Sprintf("%d", time.Now().Unix()),
		VolumeName:    volumnName,
		Status:        model.NodeStatusPending,
		Capacity:      request.Capacity,
		UsedSpace:     0,
		LastHeartbeat: time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := nm.DbManager.RegisterNode(node); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to register node: %v", err),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    node,
	})
}

func (nm *NMHandler) GetAllNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes, err := nm.DbManager.GetAllNodes()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch nodes")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    nodes,
	})
}

// GetClusterStatsHandler retrieves cluster statistics such as total nodes, total capacity,
// used capacity, and free capacity. It returns an APIResponse containing the statistics
// or an error message if the operation fails.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - r: *http.Request containing the request data.
//
// Return:
//   - Writes an APIResponse with the following structure:
//     {
//     "success": bool,
//     "message": string,
//     "data": {
//     "totalNodes": int,
//     "totalCapacity": int64,
//     "usedCapacity": int64,
//     "freeCapacity": int64,
//     },
//     "error": string,
//     }
//   - success: true if the operation is successful, false otherwise.
//   - message: Contains an error message if success is false.
//   - data: Contains the cluster statistics if success is true.
//   - error: Contains an error message if success is false.

func (nm *NMHandler) GetClusterStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := nm.DbManager.GetClusterStats()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch nodes")
		return
	}
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    stats,
	})
}

func chunkUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for form data
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	ext := filepath.Ext(filename)
	baseName := strings.TrimSuffix(filename, ext)

	data := map[string]interface{}{
		"filename": filename,
		"baseName": baseName,
		"ext":      ext,
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})

}
