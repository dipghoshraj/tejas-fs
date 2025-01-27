package cosmicmodel

type ClusterConfig struct {
	ID              string `json:"id" gorm:"primaryKey"`
	Name            string `json:"name"`
	Nodes           int64  `json:"nodes"`
	NodeCapacity    int64  `json:"nodeCapacity"`
	TotalCapacity   int64  `json:"totalCapacity"`
	UsedCapacity    int64  `json:"usedCapacity"`
	IngressNodes    int64  `json:"ingressNodes"`
	IngressCapacity int64  `json:"ingressCapacity"`
	AutoScaling     bool   `json:"autoScaling"`
}
