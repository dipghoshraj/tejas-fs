package handler

type NodeStatus string

const (
	NodeStatusActive  NodeStatus = "active"
	NodeStatusFailed  NodeStatus = "failed"
	NodeStatusPending NodeStatus = "pending"
)
