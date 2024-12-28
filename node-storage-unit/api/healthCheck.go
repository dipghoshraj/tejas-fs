package api

import (
	"encoding/json"
	"net/http"
)

type HealthCheckResponse struct {
	Status string `json:"status"` // Changed to int to return the size in MB
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Status: "active",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
