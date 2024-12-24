package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type Node struct {
	ID           string `json:"id"`
	StorageUsed  int64  `json:"storage_used"`  // in MB
	StorageLimit int64  `json:"storage_limit"` // in MB
	StoragePath  string `json:"-"`
	StorageMutex sync.Mutex
}

var node Node

func main() {
	node = Node{
		ID:           os.Getenv("NODE_ID"),
		StorageUsed:  0,
		StorageLimit: parseEnvInt("STORAGE_CAPACITY", 1024), // Default to 1024 MB
		StoragePath:  "/data",
	}

	http.HandleFunc("/status", handleStatus)
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/download", handleDownload)

	log.Printf("Node %s running with storage limit %d MB at %s\n", node.ID, node.StorageLimit, node.StoragePath)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(node)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	node.StorageMutex.Lock()
	defer node.StorageMutex.Unlock()

	// Check storage capacity
	if node.StorageUsed >= node.StorageLimit {
		http.Error(w, "Storage full", http.StatusInsufficientStorage)
		return
	}

	// Save the file
	targetPath := fmt.Sprintf("%s/%s", node.StoragePath, header.Filename)
	outFile, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	bytesWritten, err := outFile.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	node.StorageUsed += bytesWritten / (1024 * 1024) // Convert bytes to MB
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Missing file parameter", http.StatusBadRequest)
		return
	}

	targetPath := fmt.Sprintf("%s/%s", node.StoragePath, filename)
	http.ServeFile(w, r, targetPath)
}

func parseEnvInt(env string, defaultValue int64) int64 {
	val := os.Getenv(env)
	if val == "" {
		return defaultValue
	}
	var intValue int64
	_, err := fmt.Sscan(val, &intValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}
