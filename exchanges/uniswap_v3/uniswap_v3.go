package uniswapv3

import (
	"context"
	univ3factory "dex-feed/bindings/uniswap_v3"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type UniswapV3 struct {
	Client *ethclient.Client

	Factory *univ3factory.Univ3factoryCaller
}

func (v3 *UniswapV3) GetPoolWithFee(token0, token1 common.Address, fee int64) (common.Address, error) {
	// Make sure the address order is the same as in the Pool contract
	if token1.String() < token0.String() {
		tmp := token0
		token0 = token1
		token1 = tmp
	}

	// TODO: check if contract is in memory cache / in the DB

	zeroAddress := [20]byte{byte(0)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pool, err := v3.Factory.GetPool(&bind.CallOpts{Context: ctx}, token0, token1, big.NewInt(fee))
	if err != nil {
		return zeroAddress, err
	}

	if pool == zeroAddress {
		return zeroAddress, errors.New("uniswap v3 error: pool not found")
	}

	// TODO: if we get here, we should write the contract to the DB and cache it in memory

	return pool, nil
}
