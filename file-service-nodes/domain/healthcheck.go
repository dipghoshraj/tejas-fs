package domain

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/dipghoshraj/media-service/file-service-nodes/domain/proto"
)

func convertSizeToMB(size string) (int64, error) {
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
		return int64(value / 1024), nil
	case "M":
		// Already in MB
		return int64(value), nil
	default:
		// If no unit, assume the value is in MB
		return int64(value), nil
	}
}

func GetUsedSpace() (int64, error) {

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

func (s *StorageServer) HealthCheck(ctx context.Context, req *proto.HealthCheckMessage) (*proto.HelloReply, error) {
	usedSapce, err := GetUsedSpace()
	if err != nil {
		return &proto.HelloReply{
			Status:    "inactive",
			UsedSpace: 0,
		}, nil
	}
	return &proto.HelloReply{
		Status:    "active",
		UsedSpace: usedSapce,
	}, nil
}
