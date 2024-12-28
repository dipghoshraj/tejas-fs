package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type HealthCheckResponse struct {
	Status    string `json:"status"`
	UsedSapce int    `json:"usedSapce"`
	// Changed to int to return the size in MB
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func GetUsedSpace() (int, error) {

	var out bytes.Buffer
	volumePath := "/data"

	cmd := exec.Command("du", "-sh", volumePath)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("error executing du command: %v", err)
	}

	duOutput := out.String()
	duParts := strings.Fields(duOutput)

	if len(duParts) < 1 {
		return 0, fmt.Errorf("invalid output from du command")
	}

	// Parse the size and convert it to MB
	size := duParts[0]
	sizeInMB, err := convertSizeToMB(size)
	if err != nil {
		return 0, fmt.Errorf("error converting size to MB: %v", err)
	}

	// Return the size in MB
	return sizeInMB, nil
}

func convertSizeToMB(size string) (int, error) {
	// Define regex to match size values with units (e.g., 1.5G, 500M, 1024K)
	re := regexp.MustCompile(`([0-9.]+)([KMG]?)`)
	matches := re.FindStringSubmatch(size)

	if len(matches) < 3 {
		return 0, fmt.Errorf("unable to parse size: %s", size)
	}

	// Extract the numeric value and the unit
	valueStr := matches[1]
	unit := matches[2]

	// Convert the numeric value to float
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing value: %v", err)
	}

	// Convert the value to MB based on the unit
	switch unit {
	case "K":
		// KB to MB
		return int(value / 1024), nil
	case "M":
		// Already in MB
		return int(value), nil
	default:
		// If no unit, assume the value is in MB
		return int(value), nil
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	usedSapce, err := GetUsedSpace()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to decode request: %v", err),
		})
	}
	response := HealthCheckResponse{
		Status:    "active",
		UsedSapce: usedSapce,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}
