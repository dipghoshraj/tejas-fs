package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// volume stroage directory
const StorageDir = "/data"

func SaveChunk(w http.ResponseWriter, r *http.Request) {

	chunkID := r.URL.Query().Get("chunk_id")
	if chunkID == "" {
		http.Error(w, "chunk_id is required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(StorageDir, chunkID)
	file, err := os.Create(filePath)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Chunk %s saved successfully", chunkID)
}
