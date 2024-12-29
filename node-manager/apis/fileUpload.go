package apis

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/dipghoshraj/media-service/node-manager/model"
	"github.com/google/uuid"
)

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
	err = nh.DbManager.StoreEntryFile(file, dataObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dataObj)
}
