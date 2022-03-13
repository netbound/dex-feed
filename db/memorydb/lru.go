package memorydb

import lru "github.com/hashicorp/golang-lru"

// Implements Cacher
type LruCache struct {
	cache *lru.Cache
}

func New(size int) *LruCache {
	c, _ := lru.New(size)

	return &LruCache{c}
}

func (c LruCache) Get(key string) (string, bool) {
	if val, ok := c.cache.Get(key); ok {
		return val.(string), true
	}

	return "", false
}

func (c *LruCache) Put(key string, value string) {
	c.cache.Add(key, value)
}
