package handlers

import (
	"net/http"
	"HLD-REDIS-ASSIGNMENT/internal/cache"
	"HLD-REDIS-ASSIGNMENT/internal/models"
	"github.com/gin-gonic/gin"
)

type CacheHandler struct {
	cache *cache.Cache
}

func NewCacheHandler(c *cache.Cache) *CacheHandler {
	return &CacheHandler{cache: c}
}

func (h *CacheHandler) SetHandler(c *gin.Context) {
	var req models.CacheSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.cache.Set(req.Key, req.Value)
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"message": "Key inserted/updated successfully."})
}

func (h *CacheHandler) GetHandler(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ERROR",
			"message": "Key not found."})
		return
	}
	if(len(key) > 256) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "ERROR",
			"message": "Key length exceeds 256 characters."})
		return
	}
	value, exists := h.cache.Get(key)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "ERROR",
			"message": "Key not found."})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
		"key": key,
		"value": value})
}

func (h *CacheHandler) AddNodeHandler(c *gin.Context) {
	var req models.NodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.cache.AddNode(req.Node)
	c.JSON(http.StatusOK, gin.H{"message": "Node added successfully"})
}

func (h *CacheHandler) RemoveNodeHandler(c *gin.Context) {
	var req models.NodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.cache.RemoveNode(req.Node)
	c.JSON(http.StatusOK, gin.H{"message": "Node removed successfully"})
}
