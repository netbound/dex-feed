package token

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/netbound/dex-feed/db"
)

var (
	wethAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	usdcAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
)

func TestGetToken(t *testing.T) {
	rpcApi := "http://localhost:8080/eth"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, err := ethclient.DialContext(ctx, rpcApi)
	if err != nil {
		log.Fatal(err)
	}

	tdb := NewTokenDB(c, db.Opts{
		CacheSize:  20,
		Persistent: true,
	})

	token, err := tdb.GetToken(ctx, usdcAddress)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", token)
}
