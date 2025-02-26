package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const VOLUME = "/data"

func WriteSegments(fileID string, sequence int, data []byte) (string, error) {
	chunkDir := filepath.Join(VOLUME, fileID)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return "", err
	}

	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("chunk_%d", sequence))
	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		return "", err
	}

	return chunkPath, nil
}

func ReadSegment(chunkPath string) ([]byte, error) {
	data, err := os.ReadFile(chunkPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TODO : remove segments better and secure implementation
func RemoveSegment(chunkPath string) error {
	return os.Remove(chunkPath)
}
