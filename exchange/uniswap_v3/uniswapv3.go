package uniswapv3

import (
	"context"
	"errors"
	"log"
	"math/big"
	"time"

	univ3factory "dex-feed/bindings/uniswap_v3/factory"
	univ3pool "dex-feed/bindings/uniswap_v3/pool"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
)

var (
	ErrPoolNotFound = errors.New("uniswap v3 factory: pool not found")
)

type UniswapV3Addresses struct {
	FactoryAddress common.Address
}

type UniswapV3 struct {
	Client *ethclient.Client

	// Holds the pool addresses for different assets and fee tiers
	PoolCache *lru.Cache

	Factory *univ3factory.Univ3factoryCaller
	Pool    *univ3pool.Univ3poolCaller
}

// TODO
type PoolCache struct {
	MemoryCache *lru.Cache
	DbCache     interface{}
}

func New(client *ethclient.Client, addrs UniswapV3Addresses) *UniswapV3 {
	// Only errors when cache size is negative
	c, _ := lru.New(2048)

	factory, err := univ3factory.NewUniv3factoryCaller(addrs.FactoryAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	return &UniswapV3{
		Client:    client,
		PoolCache: c,
		Factory:   factory,
	}
}

func (v3 *UniswapV3) GetPoolAddress(token0, token1 common.Address, fee int64) (common.Address, error) {
	// Make sure the address order is the same as in the Pool contract, easier for lookups
	if token1.String() < token0.String() {
		tmp := token0
		token0 = token1
		token1 = tmp
	}

	// Our key are the 2 token addresses sorted and appended
	keyBytes := append(token0.Bytes(), token1.Bytes()...)
	// This works because the values are still unique
	key := string(append(keyBytes, byte(fee)))

	// TODO: check if contract is in memory cache / in the DB
	if pool, ok := v3.GetPoolCached(key); ok {
		return pool, nil
	}

	zeroAddress := [20]byte{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := v3.Factory.GetPool(&bind.CallOpts{Context: ctx}, token0, token1, big.NewInt(fee))
	if err != nil {
		return zeroAddress, err
	}

	if pool == zeroAddress {
		return zeroAddress, ErrPoolNotFound
	}

	// Add pool to memory cache
	v3.PoolCache.Add(key, pool)

	// TODO: if we get here, we should write the contract to the DB and cache it in memory

	return pool, nil
}

func (v3 *UniswapV3) GetPoolCached(key string) (common.Address, bool) {
	// TODO: check in memory, then in DB
	// if cache hit, return address and true, else return false
	if pool, ok := v3.PoolCache.Get(key); ok {
		return pool.(common.Address), true
	}

	return [20]byte{}, false
}
