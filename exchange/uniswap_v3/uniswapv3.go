package uniswapv3

import (
	"context"
	"errors"
	"log"
	"math/big"
	"time"

	univ3factory "dex-feed/bindings/uniswap_v3/factory"
	"dex-feed/db"
	"dex-feed/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrPoolNotFound = errors.New("uniswap v3 factory: pool not found")
)

type UniswapV3Addresses struct {
	FactoryAddress common.Address
}

type UniswapV3 struct {
	Client *ethclient.Client

	PoolAddressCache db.Cacher // Holds the pool addresses for different assets and fee tiers
	PoolCache        db.Cacher // Holds the actual pools in a chained cache (checks memory first, then leveldb on disk)

	Factory *univ3factory.Univ3factoryCaller
}

func New(client *ethclient.Client, addrs UniswapV3Addresses) *UniswapV3 {
	ac := db.NewCache("univ3_address_cache", 2048)
	pc := db.NewCache("univ3_pool_cache", 2048)

	factory, err := univ3factory.NewUniv3factoryCaller(addrs.FactoryAddress, client)
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

func (v3 *UniswapV3) GetPrice(token0, token1 common.Address, fee int64) (float64, error) {
	pool, err := v3.GetPool(token0, token1, fee)
	if err != nil {
		return 0, err
	}

	return pool.PriceOf(token1)
}

func (v3 *UniswapV3) GetPoolAddress(token0, token1 common.Address, fee int64) (common.Address, error) {
	// Make sure the address order is the same as in the Pool contract, easier for lookups
	token0, token1 = sortTokens(token0, token1)
	key := createPoolKey(token0, token1, fee)

	if pool, ok := v3.getPoolAddressCached(key); ok {
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

func (v3 *UniswapV3) GetPool(token0, token1 common.Address, fee int64) (*Pool, error) {
	// Make sure the address order is the same as in the Pool contract, easier for lookups
	token0, token1 = sortTokens(token0, token1)

	poolAddr, err := v3.GetPoolAddress(token0, token1, fee)
	if err != nil {
		return nil, err
	}

	key := createPoolKey(token0, token1, fee)
	if pool, ok := v3.getPoolCached(key); ok {
		return pool, nil
	}

	immutables := PoolImmutables{
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
	// PROBLEM: this does not update the state in the cache
	err = pool.UpdateState(v3.Client)
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

// UpdateCachedPoolStates should get called once the chain state updates, i.e. on a new block.
// It retrieves all the pools from the cache, updates their states and writes them to cache again.
func (v3 *UniswapV3) UpdateCachedPoolStates() error {
	iter := v3.PoolCache.NewIterator()

	for iter.Next() {
		key := iter.Key()
		poolBytes := iter.Value()
		pool, err := DecodePool(poolBytes)
		if err != nil {
			return err
		}

		err = pool.UpdateState(v3.Client)
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

func createPoolKey(token0, token1 common.Address, fee int64) string {
	// Our key is just appending the bytes of token0, token1 and the fee
	keyBytes := append(token0.Bytes(), token1.Bytes()...)
	// This works because the values are still unique
	return string(append(keyBytes, byte(fee)))
}

func sortTokens(tokenA, tokenB common.Address) (token0, token1 common.Address) {
	token0 = tokenA
	token1 = tokenB
	if tokenB.String() < tokenA.String() {
		token0 = tokenB
		token1 = tokenA
	}
	return
}
