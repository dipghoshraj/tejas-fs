package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dipghoshraj/media-service/node-manager/model"
)

type HealthResponse struct {
	Status    string `json:"status"`
	UsedSpace int    `json:"usedSapce"`
}

func FetchStatusRequest(port string) (HealthResponse, error) {
	apiURL := fmt.Sprintf("http://localhost:%s/health", port)

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error calling the API: %v", err)
		return HealthResponse{}, err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned non-OK status: %d", resp.StatusCode)
		return HealthResponse{}, fmt.Errorf("api returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the JSON response into the struct
	var health HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return HealthResponse{}, fmt.Errorf("json decoding the response: %v", err)
	}

	log.Printf("API returned health: %v", health)

	return health, nil
}

func (ndb *DBHandler) UpdateNode(response HealthResponse, node model.Node) error {
	// result := db.Model(&User{}).Where("id = ?", userID).Updates(User{Name: newName, Email: newEmail})
	log.Printf("updating %s with status %s space %d", node.ID, response.Status, response.UsedSpace)
	result := ndb.DbManager.DB.Model(&model.Node{}).Where("id = ?", node.ID).Updates(model.Node{Status: model.NodeStatus(response.Status), UsedSpace: int64(response.UsedSpace)})
	if result.Error != nil {
		log.Printf("Error updating node %v", result.Error)
		return result.Error
	}
	fmt.Printf("Rows affected: %v\n", result.RowsAffected)
	return nil
}

func (ndb *DBHandler) NodeHealth() {

	log.Printf("Node healthcheck trigger")

	nodes, err := ndb.GetAllNodes()
	if err != nil {
		log.Printf("can not get the nodes %v", err)
	}

	for _, node := range nodes {
		log.Printf("Healt check start for %s", node.ID)
		resp, err := FetchStatusRequest(node.Port)
		if err != nil {
			resp.Status = "inactive"
			log.Printf("failed to get node status %s err %v", node.ID, err)
		}

		// standardisizing the usedSpace requirement
		if resp.UsedSpace < 1 {
			resp.UsedSpace = 1
		}
		err = ndb.UpdateNode(resp, node)
		if err != nil {
			log.Printf("failed to update node status %s err %v", node.ID, err)
		}
	}

	log.Printf("Node healthcheck Ends")
}

// func main() {
// 	// Create a ticker that triggers every 10 seconds
// 	ticker := time.NewTicker(10 * time.Second)
// 	defer ticker.Stop()

// 	// Start an infinite loop that calls the API every time the ticker ticks
// 	for {
// 		select {
// 		case <-ticker.C:
// 			// Call the API each time the ticker triggers
// 			callAPI()
// 		}
// 	}
// }
