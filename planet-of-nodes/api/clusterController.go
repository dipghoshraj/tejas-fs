package api

import (
	"encoding/json"
	"net/http"

	cosmicmodel "planet-of-node/cosmic-model"

	"github.com/google/uuid"
)

type ClusterRequest struct {
	Nodes       int64 `json:"nodes"`
	NodeCpacity int64 `json:"node_capacity"`
	IngressNode int64 `json:"ingress_node"`
	AutoScale   bool  `json:"auto_scale"`
}

type Response struct {
	clsuter interface{}
}

func (na *NApi) CreateCluster(w http.ResponseWriter, r *http.Request) {

	var request ClusterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cluster_id := uuid.New().String()
	cluster := &cosmicmodel.ClusterConsfig{
		ID:           cluster_id,
		Name:         "cluster-" + cluster_id,
		Nodes:        request.Nodes,
		NodeCapacity: request.NodeCpacity,
		AutoScaling:  request.AutoScale,
		IngressNodes: request.IngressNode,
	}

	// TODO: Save the cluster metadata and create the cluster in backgroud
	if err := na.nhm.CreateCusterMetadata(cluster); err != nil {
		respondError(w, http.StatusBadRequest, "Cluster creation failed")
		return
	}

	respondWithJson(w, http.StatusOK, APIResponse{
		Success: true,
		Data: Response{
			clsuter: cluster,
		},
	})
}
