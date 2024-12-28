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

func (ndb *DBHandler) PushEntryNode(file io.Reader, filename, fileID string) error {

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
	req.Header.Set("File-ID", fileID)
	req.Header.Set("Filename", filename)

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
