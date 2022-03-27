package uniswapv3

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var uni *UniswapV3

var (
	wethAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	usdcAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	fee         = int64(500)
)

func newConnectedUniV3() *UniswapV3 {
	var (
		factoryAddress = common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")
		rpcApi         = "http://localhost:8080/eth"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, err := ethclient.DialContext(ctx, rpcApi)
	if err != nil {
		log.Fatal(err)
	}

	uni = New(c, factoryAddress, Opts{DbCache: true})
	uni.UpdateCachedPoolStates(ctx)

	return uni
}

func TestGetPoolWithFee(t *testing.T) {
	uni := newConnectedUniV3()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	t1 := time.Now()
	pool, err := uni.GetPoolAddress(ctx, wethAddress, usdcAddress, fee)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pool)

	timeNoCache := time.Since(t1)
	t.Logf("time without cache: %v", timeNoCache)

	// if pool != common.HexToAddress("0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8") {
	// 	t.Fatalf("wrong pool: %s", pool)
	// }

	t2 := time.Now()
	pool, err = uni.GetPoolAddress(ctx, wethAddress, usdcAddress, fee)
	if err != nil {
		t.Fatal(err)
	}

	timeWithCache := time.Since(t2)
	t.Logf("time with cache: %v", timeWithCache)
}

func TestGetPrice(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	price, err := uni.GetPrice(ctx, wethAddress, usdcAddress, fee)
	if err != nil {
		t.Log(err)
	}

	t.Log(wethAddress, usdcAddress, price)

	price, err = uni.GetPrice(ctx, usdcAddress, wethAddress, fee)
	if err != nil {
		t.Log(err)
	}

	t.Log(usdcAddress, wethAddress, price)
}
