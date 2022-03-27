package uniswapv3

import (
	"context"
	"errors"
	"log"
	"math/big"

	"github.com/netbound/dex-feed/db"
	"github.com/netbound/dex-feed/db/memorydb"
	"github.com/netbound/dex-feed/token"

	univ3factory "github.com/netbound/dex-feed/bindings/uniswap_v3/factory"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrPoolNotFound = errors.New("uniswap v3 factory: pool not found")
)

type UniswapV3 struct {
	Client *ethclient.Client

	PoolAddressCache db.Cacher // Holds the pool addresses for different assets and fee tiers
	PoolCache        db.Cacher // Holds the actual pools in a chained cache (checks memory first, then leveldb on disk)

	TokenCache token.TokenCache

	Factory *univ3factory.Univ3factoryCaller

	Opts Opts // Holds options
}

type Opts struct {
	DbCache   bool // Should cache be persisted? (leveldb)
	CacheSize int  // Size of the in-memory LRU cache
}

// Returns a UniswapV3 instance.
func New(client *ethclient.Client, factoryAddress common.Address, opts Opts) *UniswapV3 {
	var ac, pc db.Cacher

	cacheSize := opts.CacheSize
	if cacheSize == 0 {
		cacheSize = 2048
	}

	// By default, only use memory cache
	ac = memorydb.New(cacheSize)
	pc = memorydb.New(cacheSize)

	// If dbCache flag is set, initalize leveldb
	if opts.DbCache {
		ac = db.NewDBCache("univ3_address_cache", cacheSize)
		pc = db.NewDBCache("univ3_pool_cache", cacheSize)
	}

	factory, err := univ3factory.NewUniv3factoryCaller(factoryAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	return &UniswapV3{
		Client:           client,
		PoolAddressCache: ac,
		PoolCache:        pc,
		Factory:          factory,
	}
}

func (v3 *UniswapV3) GetPrice(ctx context.Context, token0, token1 common.Address, fee int64) (float64, error) {
	pool, err := v3.GetPool(ctx, token0, token1, fee)
	if err != nil {
		return 0, err
	}

	return pool.PriceOf(token1)
}

func (v3 *UniswapV3) GetPoolAddress(ctx context.Context, token0, token1 common.Address, fee int64) (common.Address, error) {
	// Make sure the address order is the same as in the Pool contract, easier for lookups
	token0, token1 = sortTokens(token0, token1)
	key := createPoolKey(token0, token1, fee)

	if pool, ok := v3.getPoolAddressCached(key); ok {
		return pool, nil
	}

	zeroAddress := [20]byte{}

	pool, err := v3.Factory.GetPool(&bind.CallOpts{Context: ctx}, token0, token1, big.NewInt(fee))
	if err != nil {
		return zeroAddress, err
	}

	if pool == zeroAddress {
		return zeroAddress, ErrPoolNotFound
	}

	// Cache the pool (both in-memory and on-disk)
	v3.PoolAddressCache.Put(key, pool.Bytes())

	return pool, nil
}

func (v3 *UniswapV3) getPoolAddressCached(key string) (common.Address, bool) {
	// if cache hit, return address and true, else return false
	if pool, ok := v3.PoolAddressCache.Get(key); ok {
		return common.BytesToAddress(pool), true
	}

	return [20]byte{}, false
}

func (v3 *UniswapV3) GetPool(ctx context.Context, token0, token1 common.Address, fee int64) (*Pool, error) {
	// Make sure the address order is the same as in the Pool contract, easier for lookups
	token0, token1 = sortTokens(token0, token1)

	poolAddr, err := v3.GetPoolAddress(ctx, token0, token1, fee)
	if err != nil {
		return nil, err
	}

	key := createPoolKey(token0, token1, fee)
	if pool, ok := v3.getPoolCached(key); ok {
		return pool, nil
	}

	immutables := PoolOpts{
		// TODO: hardcoded decimals for now
		Token0: token.Token{Address: token0, Decimals: 6, Name: "USDC"},
		Token1: token.Token{Address: token1, Decimals: 18, Name: "WETH"},
		Fee:    fee,
	}

	pool, err := NewPool(v3.Client, "USDCWETH", poolAddr, immutables)
	if err != nil {
		return nil, err
	}

	// First time: update initial state
	err = pool.UpdateState(ctx, v3.Client)
	if err != nil {
		return nil, err
	}

	encoded, err := pool.EncodePool()
	if err != nil {
		return nil, err
	}

	v3.PoolCache.Put(key, encoded)

	return pool, nil
}

func (v3 *UniswapV3) getPoolCached(key string) (*Pool, bool) {
	if poolBytes, ok := v3.PoolCache.Get(key); ok {
		// TODO: convert to pool
		pool, err := DecodePool(poolBytes)
		if err != nil {
			return nil, false
		}
		return pool, true
	}

	return nil, false
}

// TODO
func (v3 *UniswapV3) GetToken(address common.Address) (token.Token, error) {

	return token.Token{}, nil
}

// UpdateCachedPoolStates should get called once the chain state updates, i.e. on a new block.
// It retrieves all the pools from the cache, updates their states and writes them to cache again.
func (v3 *UniswapV3) UpdateCachedPoolStates(ctx context.Context) error {
	iter := v3.PoolCache.NewIterator()

	for iter.Next() {
		key := iter.Key()
		poolBytes := iter.Value()
		pool, err := DecodePool(poolBytes)
		if err != nil {
			return err
		}

		err = pool.UpdateState(ctx, v3.Client)
		if err != nil {
			return err
		}

		encoded, err := pool.EncodePool()
		if err != nil {
			return err
		}

		v3.PoolCache.Put(string(key), encoded)
	}

	iter.Release()
	return nil
}
