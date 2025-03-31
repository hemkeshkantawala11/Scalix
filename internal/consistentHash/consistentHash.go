package consistentHash

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
}

func New(replicas int, fn Hash) *ConsistentHash {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
}

func (c *ConsistentHash) Add(nodes ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, node := range nodes {
		for i := 0; i < c.replicas; i++ {
			hash := int(c.hash([]byte(node + string(i))))
			c.keys = append(c.keys, hash)
			c.hashMap[hash] = node
		}
	}
	sort.Ints(c.keys)
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
