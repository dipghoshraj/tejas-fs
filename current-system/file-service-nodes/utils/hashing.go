package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func generateSegmentsHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func IntegrityCheck(data []byte, hashDigests string) error {
	calculated := generateSegmentsHash(data)
	if calculated != hashDigests {
		return fmt.Errorf("segments HashDigests is not valid")
	}
	return nil
}

func FileHashDigests(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file path: %v", err)
	}

	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil

}
