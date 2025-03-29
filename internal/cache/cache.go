package cache

import (
	"container/heap"
	"sync"
	"time"
)

// CacheItem represents a single item in the cache
type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration int64
	Accessed   int64 // For LRU tracking
	index      int   // Index for heap operations
}

// CustomCache is a thread-safe in-memory cache
type CustomCache struct {
	items      map[string]*CacheItem
	mutex      sync.RWMutex
	capacity   int
	evictionQ  priorityQueue
	evictChan  chan string
}

// New creates a new custom cache
func New(options ...Option) *CustomCache {
	cache := &CustomCache{
		items:     make(map[string]*CacheItem),
		capacity:  defaultCapacity,
		evictChan: make(chan string, 100),
		evictionQ: make(priorityQueue, 0),
	}

	heap.Init(&cache.evictionQ)

	for _, option := range options {
		option(cache)
	}

	if cache.capacity > 0 {
		go cache.startEvictionProcess()
	}

	return cache
}

const (
	defaultCapacity = 1000
	defaultTTL      = time.Hour
)

// Option represents a function that modifies CustomCache
type Option func(*CustomCache)

// WithCapacity sets the capacity of the cache
func WithCapacity(capacity int) Option {
	return func(c *CustomCache) {
		c.capacity = capacity
	}
}

// Set adds or updates an item in the cache
func (c *CustomCache) Set(key string, value interface{}, ttl ...time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	duration := defaultTTL
	if len(ttl) > 0 {
		duration = ttl[0]
	}

	expiration := time.Now().Add(duration).UnixNano()
	now := time.Now().UnixNano()

	if len(c.items) >= c.capacity {
		c.evict() // Evict only when necessary
	}

	item := &CacheItem{
		Key:        key,
		Value:      value,
		Expiration: expiration,
		Accessed:   now,
	}

	c.items[key] = item
	heap.Push(&c.evictionQ, item)
}

// Get retrieves an item from the cache
func (c *CustomCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	item, exists := c.items[key]
	c.mutex.RUnlock()

	if !exists || time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	// Update LRU
	c.mutex.Lock()
	item.Accessed = time.Now().UnixNano()
	heap.Fix(&c.evictionQ, item.index)
	c.mutex.Unlock()

	return item.Value, true
}

// evict removes the least recently used item
func (c *CustomCache) evict() {
	if len(c.items) == 0 {
		return
	}

	item := heap.Pop(&c.evictionQ).(*CacheItem)
	delete(c.items, item.Key)
}

// CleanExpired removes expired items
func (c *CustomCache) CleanExpired() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, item := range c.items {
		if time.Now().UnixNano() > item.Expiration {
			delete(c.items, key)
		}
	}
}

// startEvictionProcess runs a background process to clean expired items
func (c *CustomCache) startEvictionProcess() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		c.CleanExpired()
	}
}

// Priority queue (Min-Heap) implementation
type priorityQueue []*CacheItem

func (pq priorityQueue) Len() int           { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool { return pq[i].Accessed < pq[j].Accessed }
func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*CacheItem)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
