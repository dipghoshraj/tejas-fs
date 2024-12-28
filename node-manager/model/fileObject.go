package model

type object struct {
	id            string
	name          string
	size          int64
	totalChunks   int
	entry_node_id string
	distributed   bool
}
