package handler

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/dipghoshraj/media-service/node-manager/model"
)

func (ndb *DBHandler) FindEntryPointNode() (*model.Node, error) {

	var entrypoint *model.Node
	result := ndb.DbManager.DB.Where("status = ? and capacity - used_space >= ?", "active", 5).Order("capacity - used_space DESC").First(&entrypoint)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find entry point node: %v", result.Error)
	}
	return entrypoint, nil
}

func (ndb *DBHandler) StoreEntryFile(file io.Reader, dataObject model.DataObject) (model.DataObject, error) {

	entrypoint, err := ndb.FindEntryPointNode()
	if err != nil {
		return dataObject, fmt.Errorf("failed to find entry point node: %v", err)
	}

	dataObject.EntryNodeId = entrypoint.ID
	tx := ndb.DbManager.DB.Begin()
	fmt.Println(dataObject)

	if err := tx.Create(dataObject).Error; err != nil {
		tx.Rollback()
		return dataObject, fmt.Errorf("failed to create entry point %v", err)
	}

	err = ndb.RequestStore(file, dataObject.Ext, dataObject.ID)
	if err != nil {
		tx.Rollback()
		return dataObject, fmt.Errorf("failed to push entry point %v", err)
	}

	return dataObject, tx.Commit().Error
}

func (ndb *DBHandler) RequestStore(file io.Reader, ext string, fileId string) error {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("file_id", fileId)
	_ = writer.WriteField("ext", ext)

	part3, errFile3 := writer.CreateFormFile("datafile", fileId)
	if errFile3 != nil {
		return fmt.Errorf("failed to create file payload %v", errFile3)
	}
	_, errFile3 = io.Copy(part3, file)

	if errFile3 != nil {
		return fmt.Errorf("failed to set file payload %v", errFile3)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer %v", err)
	}

	entrypoint, err := ndb.FindEntryPointNode()
	if err != nil {
		return fmt.Errorf("failed to find entry point node: %v", err)
	}

	// const entryNodeURL = "http://entry-node:8081/store" // Replace with your Entry Node URL
	entryNodeURL := fmt.Sprintf("http://localhost:%s/store", entrypoint.Port)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("POST", entryNodeURL, payload)

	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to %v", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to upload file to Entry Node, status: %d", resp.StatusCode)
	}
	return nil
}
