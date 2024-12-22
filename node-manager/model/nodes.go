package model

type Node struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
	LastPing int64  `json:"last_ping"`
}
