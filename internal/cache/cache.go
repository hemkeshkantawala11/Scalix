package cache

import (
    "log"
    "runtime"
    "time"
    "sync"
    "github.com/hashicorp/golang-lru"
    consistenthash "HLD-REDIS-ASSIGNMENT/internal/consistentHash"
)

type Cache struct {
    shards   map[string]*sync.Map
    setCh    chan [6]string
    hash     *consistenthash.ConsistentHash
    mu       sync.RWMutex
    lruCache *lru.Cache // LRU Cache
}


func New(nodes []string, cacheSize int) *Cache {
    lru, err := lru.New(cacheSize) // Create LRU cache with fixed size
    if err != nil {
        log.Fatalf("Error creating LRU cache: %v", err)
    }

    cache := &Cache{
        shards:   make(map[string]*sync.Map),
        setCh:    make(chan [6]string, 50000),
        hash:     consistenthash.New(100, nil),
        lruCache: lru, // Assign LRU cache
    }

    for _, node := range nodes {
        cache.shards[node] = &sync.Map{}
        cache.hash.Add(node)
    }

    // Start background memory monitoring
    go cache.monitorMemoryUsage()

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
    c.mu.Lock()
    defer c.mu.Unlock()

    // Store in LRU cache
    c.lruCache.Add(key, value)

    // Add to set channel
    select {
    case c.setCh <- [6]string{key, value}:
    default:
        log.Println("Set channel is full, dropping request")
    }
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    if value, ok := c.lruCache.Get(key); ok {
        return value.(string), true
    }

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

func (c *Cache) monitorMemoryUsage() {
    for {
        var memStats runtime.MemStats
        runtime.ReadMemStats(&memStats)

        memUsage := float64(memStats.Alloc) / float64(memStats.Sys) * 100
        if memUsage > 70 {
            c.mu.Lock()
            c.lruCache.Purge() // Evict all items
            log.Println("Eviction triggered due to high memory usage (>70%)")
            c.mu.Unlock()
        }
        time.Sleep(5 * time.Second) // Check memory usage every 5 seconds
    }
}

