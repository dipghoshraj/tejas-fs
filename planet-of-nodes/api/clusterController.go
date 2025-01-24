package api

import "net/http"

type Cluster struct {
	Nodes        int64 `json:"nodes"`
	NodecCpacity int64 `json:"node_capacity"`
	IngressNode  int64 `json:"ingress_node"`
	AutoScale    bool  `json:"auto_scale"`
}

type Response struct{}

func (na *NApi) CreateCluster(w http.ResponseWriter, r *http.Request) {

	respondWithJson(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    Response{},
	})

}
