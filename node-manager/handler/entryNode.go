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

func (ndb *DBHandler) FindIngressNode() (*model.Node, error) {

	var entrypoint *model.Node
	result := ndb.DbManager.DB.Where("status = ? and capacity - used_space >= ?", "active", 5).Order("capacity - used_space DESC").First(&entrypoint)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to find entry point node: %v", result.Error)
	}
	return entrypoint, nil
}

func (ndb *DBHandler) StoreEntryFile(file io.Reader, orb model.Orbs) (model.Orbs, error) {

	ingressNode, err := ndb.FindIngressNode()
	if err != nil {
		return orb, fmt.Errorf("failed to find entry point node: %v", err)
	}

	orb.IngressNodeId = ingressNode.ID
	tx := ndb.DbManager.DB.Begin()
	fmt.Println(orb)

	if err := tx.Create(orb).Error; err != nil {
		tx.Rollback()
		return orb, fmt.Errorf("failed to create entry point %v", err)
	}

	err = ndb.RequestStore(file, orb.Ext, orb.ID, ingressNode.Port)
	if err != nil {
		tx.Rollback()
		return orb, fmt.Errorf("failed to push entry point %v", err)
	}

	return orb, tx.Commit().Error
}

func (ndb *DBHandler) RequestStore(file io.Reader, ext string, fileId string, ingressPort string) error {
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

	entryNodeURL := fmt.Sprintf("http://localhost:%s/store", ingressPort)
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
