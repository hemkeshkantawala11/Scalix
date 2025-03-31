package consistenthash

import (
	"hash/crc32"
	"sort"
	"sync"
)

type Hash func(data []byte) uint32

type ConsistentHash struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
	mu       sync.RWMutex
	nodes    map[string]bool
}

func New(replicas int, fn Hash) *ConsistentHash {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
		nodes:    make(map[string]bool),
	}
}

func (c *ConsistentHash) Add(nodes ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, node := range nodes {
		if c.nodes[node] {
			continue // Node already exists
		}
		c.nodes[node] = true

		for i := 0; i < c.replicas; i++ {
			hash := int(c.hash([]byte(node + string(i))))
			c.keys = append(c.keys, hash)
			c.hashMap[hash] = node
		}
	}
	sort.Ints(c.keys)
}

func (c *ConsistentHash) Remove(node string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.nodes[node] {
		return // Node not present
	}
	delete(c.nodes, node)

	newKeys := []int{}
	newHashMap := make(map[int]string)

	for _, key := range c.keys {
		if c.hashMap[key] != node {
			newKeys = append(newKeys, key)
			newHashMap[key] = c.hashMap[key]
		}
	}

	c.keys = newKeys
	c.hashMap = newHashMap
}

func (c *ConsistentHash) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.keys) == 0 {
		return ""
	}

	hash := int(c.hash([]byte(key)))
	idx := sort.Search(len(c.keys), func(i int) bool { return c.keys[i] >= hash })

	if idx == len(c.keys) {
		idx = 0
	}

	return c.hashMap[c.keys[idx]]
}
