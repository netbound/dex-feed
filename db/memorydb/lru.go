package memorydb

import (
	lru "github.com/hashicorp/golang-lru"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

// Implements Cacher
type LruCache struct {
	cache *lru.Cache
}

func New(size int) *LruCache {
	c, _ := lru.New(size)

	return &LruCache{c}
}

func (c LruCache) Get(key string) ([]byte, bool) {
	if val, ok := c.cache.Get(key); ok {
		return val.([]byte), true
	}

	return []byte{}, false
}

func (c *LruCache) Put(key string, value []byte) {
	c.cache.Add(key, value)
}

// Hack, db.Cacher interface needs this
func (c LruCache) NewIterator() iterator.Iterator {
	return iterator.NewEmptyIterator(nil)
}
