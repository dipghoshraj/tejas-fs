package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const segmentSize = 4 * 1024

func WriteSegmentFile(filePath string, outputDir string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	segments := []string{}
	buffer := make([]byte, segmentSize)
	sequence := 0

	for {
		n, err := file.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}
		if n == 0 {
			break
		}

		segmentPath := filepath.Join(outputDir, fmt.Sprintf("%s_segment_%d", filepath.Base(filePath), sequence))
		if err := os.WriteFile(segmentPath, buffer[:n], 0644); err != nil {
			return nil, err
		}

		segments = append(segments, segmentPath)
		sequence++
	}

	return segments, nil
}
