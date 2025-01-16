package model

type Orbs struct {
	ID            string
	Name          string
	Size          int64
	TotalChunks   int
	IngressNodeId string
	Distributed   bool
	Ext           string
}
