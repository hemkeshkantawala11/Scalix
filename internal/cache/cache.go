package cache

import (
	"sync"
)

type Cache struct {
	data  sync.Map
	setCh chan [2]string
}

func New() *Cache {
	cache := &Cache{
		setCh: make(chan [2]string, 1000),
	}
	go cache.processSetOperations()
	return cache
}

func (c *Cache) processSetOperations() {
	for pair := range c.setCh {
		c.data.Store(pair[0], pair[1])
	}
}

func (c *Cache) Set(key, value string) {
	c.setCh <- [2]string{key, value}
}

func (c *Cache) Get(key string) (string, bool) {
	val, exists := c.data.Load(key)
	if !exists {
		return "", false
	}
	return val.(string), true
}