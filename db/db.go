package db

import (
	"log"

	"github.com/netbound/dex-feed/db/leveldb"
	"github.com/netbound/dex-feed/db/memorydb"

	"github.com/syndtr/goleveldb/leveldb/iterator"
)

type Cacher interface {
	Get(key string) ([]byte, bool)
	Put(key string, value []byte)

	NewIterator() iterator.Iterator
}

type Cache struct {
	lruCache *memorydb.LruCache
	dbCache  *leveldb.Database
}

func NewCache(name string, size int) *Cache {
	dbCache, err := leveldb.NewDatabase(name)
	if err != nil {
		log.Fatal("err creating leveldb", err)
	}

	return &Cache{
		lruCache: memorydb.New(size),
		dbCache:  dbCache,
	}
}

func (c Cache) Get(key string) ([]byte, bool) {
	// First check in-memory cache
	if val, ok := c.lruCache.Get(key); ok {
		return val, true
	}

	// Then check on-disk cache
	if val, ok := c.dbCache.Get(key); ok {
		// If found on disk, cache in memory for later hits
		c.lruCache.Put(key, val)
		return val, true
	}

	return []byte{}, false
}

func (c *Cache) Put(key string, value []byte) {
	c.lruCache.Put(key, value)
	c.dbCache.Put(key, value)
}

func (c Cache) NewIterator() iterator.Iterator {
	return c.dbCache.NewIterator()
}
