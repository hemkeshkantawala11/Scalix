package handlers

import (
	"net/http"
	"time"

	"HLD-REDIS-ASSIGNMENT/internal/cache"
	"HLD-REDIS-ASSIGNMENT/internal/models"
	"github.com/gin-gonic/gin"
)

// CacheHandler manages HTTP interactions with the cache
type CacheHandler struct {
	cache *cache.CustomCache
}

// NewCacheHandler creates a new cache handler
func NewCacheHandler(c *cache.CustomCache) *CacheHandler {
	return &CacheHandler{cache: c}
}

// SetHandler handles cache insert/update
func (h *CacheHandler) SetHandler(c *gin.Context) {
	var req models.CacheSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "message": err.Error()})
		return
	}

	ttl := time.Hour
	if req.TTL > 0 {
		ttl = time.Duration(req.TTL) * time.Second
	}

	h.cache.Set(req.Key, req.Value, ttl)

	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Key inserted/updated"})
}

// GetHandler retrieves cache item
func (h *CacheHandler) GetHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "message": "Key is required"})
		return
	}
	if len(key) > 256 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "ERROR", "message": "Key exceeds maximum length of 256 characters"})
		return
	}

	value, found := h.cache.Get(key)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"status": "ERROR", "message": "Key not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK", "key": key, "value": value})
}

// StatsHandler returns cache statistics

