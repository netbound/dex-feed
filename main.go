package main

import (
	"context"
	"fmt"
	"log"
	"time"

	univ3factory "dex-feed/bindings/uniswap_v3/factory"
	uniswapv3 "dex-feed/exchanges/uniswap_v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	var (
		factoryAddress = common.HexToAddress("0x1F98431c8aD98523631AE4a59f267346ea31F984")

		wethAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
		usdcAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
		fee         = int64(3000)

		rpcApi = "http://localhost:8080/eth"
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err := ethclient.DialContext(ctx, rpcApi)
	if err != nil {
		log.Fatal(err)
	}

	factoryInstance, err := univ3factory.NewUniv3factoryCaller(factoryAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	univ3 := uniswapv3.UniswapV3{
		Client:  client,
		Factory: factoryInstance,
	}

	pool, err := univ3.GetPoolAddress(wethAddress, usdcAddress, fee)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pool:", pool)
}
