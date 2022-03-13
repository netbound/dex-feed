package db

import "dex-feed/db/memorydb"

type Cacher interface {
	Get(key string) (interface{}, bool)
	Put(key string, value interface{})
}

type Cache struct {
	lruCache *memorydb.LruCache
}

func NewCache(size int) *Cache {
	return &Cache{
		lruCache: memorydb.New(size),
	}
}

func (c Cache) Get(key string) (interface{}, bool) {
	if val, ok := c.lruCache.Get(key); ok {
		return val, true
	}

	return nil, false
}

func (c *Cache) Put(key string, value interface{}) {
	c.lruCache.Put(key, value)
}
