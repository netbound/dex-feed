package uniswapv3

import "github.com/ethereum/go-ethereum/common"

type PoolManager interface {
	GetPool(token0, token1 common.Address, fee int64) (Pool, error)
	GetPoolAddress(token0, token1 common.Address, fee int64) (common.Address, error)
}

type poolManager struct{}
