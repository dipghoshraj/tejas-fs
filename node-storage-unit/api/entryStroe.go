package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type StoreResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func EntryPointStoreage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit for form data
	if err != nil {
		http.Error(w, "File size too large", http.StatusForbidden)
		return
	}

	file_id := r.FormValue("file_id")
	file, handler, err := r.FormFile("datafile")

	if err != nil {
		http.Error(w, "Unable to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	originalFilename := handler.Filename
	ext := filepath.Ext(originalFilename)

	newFileName := fmt.Sprintf("%s.%s", file_id, ext)
	filePath := filepath.Join(StorageDir, newFileName)

	dst, err := os.Create(filePath)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create file: %v", err), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to save file: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(StoreResponse{
		Success: true,
		Message: "file created",
	})
}
