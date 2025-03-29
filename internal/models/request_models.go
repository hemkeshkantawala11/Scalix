package models

// CacheSetRequest represents the structure for setting a cache item
type CacheSetRequest struct {
	Key   string `json:"key" binding:"required,max=256"`
	Value string `json:"value" binding:"required,max=256"`
	TTL   int    `json:"ttl,omitempty"` // Optional TTL in seconds
}