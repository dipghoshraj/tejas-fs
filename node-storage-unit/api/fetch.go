package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func GetChunk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chunkID := vars["chunkID"]

	if chunkID == "" {
		http.Error(w, "chunk_id is required", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(StorageDir, chunkID)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve chunk: %v", err), http.StatusInternalServerError)
	}
}
