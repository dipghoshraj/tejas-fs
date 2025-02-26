package storage

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func ConverToMB(size string) (int64, error) {

	re := regexp.MustCompile(`([0-9.]+)([KMG]?)`)
	matches := re.FindStringSubmatch(size)

	if len(matches) < 3 {
		return 0, fmt.Errorf("unable to parse size: %s", size)
	}

	// Extract the numeric value and the unit
	valueStr := matches[1]
	unit := matches[2]
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing value: %v", err)
	}

	switch unit {
	case "K":
		return int64(value / 1024), nil
	case "M":
		return int64(value), nil
	default:
		return int64(value), nil
	}
}

func GetUsedSpace() (int64, error) {

	var out bytes.Buffer

	// IMPORATNT : DO NOT CHANGE THIS VARIABLE
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
	sizeInMB, err := ConverToMB(size)
	if err != nil {
		return 0, fmt.Errorf("error converting size to MB: %v", err)
	}

	// Return the size in MB
	return sizeInMB, nil
}
