package token

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/netbound/dex-feed/db"
	"github.com/netbound/dex-feed/db/memorydb"
)

type TokenCache struct {
	cache db.Cacher
}

func NewCache(size int, persistent bool) TokenCache {
	var c db.Cacher
	c = memorydb.New(size)

	if persistent {
		c = db.NewDBCache("token_cache", size)
	}

	return TokenCache{
		cache: c,
	}
}

func (tc *TokenCache) Add(token Token) {
	// We can ignore the error here
	encoded, _ := token.Encode()
	tc.cache.Put(token.Address.String(), encoded)
}

func (tc TokenCache) Get(address common.Address) (Token, bool) {
	if encoded, ok := tc.cache.Get(address.String()); ok {
		t, err := Decode(encoded)
		if err != nil {
			return Token{}, false
		}

		return t, true
	}

	return Token{}, false
}
