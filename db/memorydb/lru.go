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

func (c LruCache) Get(key string) (interface{}, bool) {
	return c.cache.Get(key)
}

func (c *LruCache) Put(key string, value interface{}) {
	c.cache.Add(key, value)
}
