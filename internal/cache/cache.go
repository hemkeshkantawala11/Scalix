package cache

import (
	"sync"
	consistenthash "HLD-REDIS-ASSIGNMENT/internal/consistentHash"
	"log"
)

type Cache struct {
	shards map[string]*sync.Map
	setCh  chan [6]string
	hash   *consistenthash.ConsistentHash
	mu     sync.RWMutex
}

func New(nodes []string) *Cache {
	cache := &Cache{
		shards: make(map[string]*sync.Map),
		setCh:  make(chan [6]string, 50000),
		hash:   consistenthash.New(100, nil), 
	}

	for _, node := range nodes {
		cache.shards[node] = &sync.Map{}
		cache.hash.Add(node)
	}

	for i := 0; i < 10; i++ { 
        go cache.processSetOperations()
    }
	return cache
}

func (c *Cache) getShard(key string) *sync.Map {
	node := c.hash.Get(key)
	return c.shards[node]
}

func (c *Cache) processSetOperations() {
	for pair := range c.setCh {
		shard := c.getShard(pair[0])
		shard.Store(pair[0], pair[1])
	}
}

func (c *Cache) Set(key, value string) {
	select {
    case c.setCh <- [6]string{key, value}:
    default:
        log.Println("Set channel is full, dropping request")
    }
}

func (c *Cache) Get(key string) (string, bool) {
	shard := c.getShard(key)
	val, exists := shard.Load(key)
	if !exists {
		return "", false
	}
	return val.(string), true
}

func (c *Cache) AddNode(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.shards[node]; !exists {
		c.shards[node] = &sync.Map{}
		c.hash.Add(node)
	}
}

func (c *Cache) RemoveNode(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.shards[node]; exists {
		delete(c.shards, node)
		c.hash.Remove(node)
	}
}
