package models

type CacheSetRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type NodeRequest struct {
	Node string `json:"node" binding:"required"`
}
