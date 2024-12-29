package handler

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dipghoshraj/media-service/node-manager/model"
)

func (ndb *DBHandler) FindEntryPointNode() (*model.Node, error) {

	var entrypoint *model.Node
	result := ndb.DbManager.DB.Where("status = ?, capacity - used_space >= ?", "active", 5).First(&entrypoint)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find entry point node: %v", result.Error)
	}
	return entrypoint, nil
}

func (ndb *DBHandler) StoreEntryFile(file io.Reader, dataObject model.DataObject) error {
	entrypoint, err := ndb.FindEntryPointNode()
	if err != nil {
		return fmt.Errorf("failed to find entry point node: %v", err)
	}

	dataObject.EntryNodeId = entrypoint.ID
	tx := ndb.DbManager.DB.Begin()

	if err := tx.Create(dataObject).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create entry point %v", err)
	}
	err = ndb.PushEntryNode(file, dataObject.Ext, dataObject.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to push entry point %v", err)
	}
	return nil
}

func (ndb *DBHandler) PushEntryNode(file io.Reader, ext string, fileID string) error {

	entrypoint, err := ndb.FindEntryPointNode()
	if err != nil {
		return fmt.Errorf("failed to find entry point node: %v", err)
	}

	// const entryNodeURL = "http://entry-node:8081/store" // Replace with your Entry Node URL
	entryNodeURL := fmt.Sprintf("http://localhost:%s/store", entrypoint.Port)
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", entryNodeURL, file)
	if err != nil {
		return err
	}
	req.Header.Set("file_id", fileID)
	req.Header.Set("ext", ext)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file to Entry Node, status: %d", resp.StatusCode)
	}
	return nil
}
