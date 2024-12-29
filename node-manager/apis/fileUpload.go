package apis

import (
	"net/http"
	"path/filepath"

	"github.com/dipghoshraj/media-service/node-manager/model"
	"github.com/google/uuid"
)

// SaveFile handles file uploads and stores the metadata in the database.
// It accepts a POST request with a form-data containing a file under the key "myFile".
// The function parses the form data, extracts the file and its metadata, generates a unique ID for the file,
// and stores the file and its metadata in the database.
//
// Parameters:
// - w: http.ResponseWriter to write the response.
// - r: *http.Request containing the form-data with the file.
//
// Returns:
// - If successful, writes a JSON response with the stored file's metadata and sets the HTTP status code to 200 (StatusOK).
// - If any error occurs during form parsing, file extraction, database storage, or response writing,
//   writes an error message and sets the HTTP status code accordingly.

type Response struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	TotalChunks int    `json:"totalChunks"`
	EntryNodeId string `json:"entryNodeId"`
	Distributed bool   `json:"distributed"`
	Ext         string `json:"ext"`
	ReplicaId   string `json:"replicaId"`
}

func (nh *NMHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	originalFilename := handler.Filename
	ext := filepath.Ext(originalFilename)

	dataObj := model.DataObject{
		ID:          uuid.New().String(),
		Ext:         ext,
		Name:        originalFilename,
		Distributed: false,
	}
	dataObj, err = nh.DbManager.StoreEntryFile(file, dataObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		ID:          dataObj.ID,
		Name:        dataObj.Name,
		EntryNodeId: dataObj.EntryNodeId,
		ReplicaId:   dataObj.ReplicaId,
		Ext:         dataObj.Ext,
		TotalChunks: dataObj.TotalChunks,
		Size:        dataObj.Size,
		Distributed: dataObj.Distributed,
	})
}
