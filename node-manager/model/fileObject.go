package model

type DataObject struct {
	ID          string
	Name        string
	Size        int64
	TotalChunks int
	EntryNodeId string
	Distributed bool
	Ext         string
	ReplicaId   string
}
