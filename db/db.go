package db

import (
	"dex-feed/db/memorydb"
	"dex-feed/db/postgresdb"
	"log"
)

type Cacher interface {
	Get(key string) (string, bool)
	Put(key string, value string)
}

type Cache struct {
	lruCache *memorydb.LruCache
	dbCache  *postgresdb.PostgresDb
}

func NewCache(name string, size int) *Cache {
	connStr := "postgres://dex-feed:dex-feed@localhost/dex-feed?sslmode=disable"
	dbCache, err := postgresdb.NewDbCache(connStr, name)
	if err != nil {
		log.Fatal(err)
	}

	return &Cache{
		lruCache: memorydb.New(size),
		dbCache:  dbCache,
	}
}

func (c Cache) Get(key string) (string, bool) {
	// First check in-memory cache
	if val, ok := c.lruCache.Get(key); ok {
		return val, true
	}

	// Then check on-disk cache
	if val, ok := c.dbCache.Get(key); ok {
		return val, true
	}

	return "", false
}

func (c *Cache) Put(key string, value string) {
	c.lruCache.Put(key, value)
	c.dbCache.Put(key, value)
}
