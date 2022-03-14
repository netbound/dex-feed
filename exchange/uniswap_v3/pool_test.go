package uniswapv3

import (
	"context"
	"dex-feed/token"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethUsdcPool = common.HexToAddress("0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8")

	wethAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	usdcAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	fee         = int64(3000)

	token1 = token.NewToken("WETH", wethAddress, 18)
	token0 = token.NewToken("USDC", usdcAddress, 6)
)

func newTestPool() (*Pool, error) {

	i := PoolImmutables{
		Token0: token0,
		Token1: token1,
		Fee:    fee,
	}

	rpcApi := "http://localhost:8080/eth"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	c, err := ethclient.DialContext(ctx, rpcApi)
	if err != nil {
		return nil, err
	}

	p, err := NewPool(c, "WETHUSDC", ethUsdcPool, i)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func TestNewPool(t *testing.T) {
	p, err := newTestPool()
	if err != nil {
		t.Fatalf("error creating pool: %s", err)
	}

	t.Log("new pool name:", p.Name)
}

func TestUpdateState(t *testing.T) {
	p, err := newTestPool()
	if err != nil {
		t.Fatalf("error creating pool: %s", err)
	}

	err = p.UpdateState()
	if err != nil {
		t.Fatalf("error updating state: %s", err)
	}

	emptyState := PoolState{}

	if p.State == emptyState {
		t.Fatal("empty state")
	}

	t.Log("pool state:", p.State)
}

func TestPriceOf(t *testing.T) {
	p, err := newTestPool()
	if err != nil {
		t.Fatalf("error creating pool: %s", err)
	}

	err = p.UpdateState()
	if err != nil {
		t.Fatalf("error updating state: %s", err)
	}

	price, err := p.PriceOf(token1)
	if err != nil {
		t.Fatalf("error getting price: %s", err)
	}

	t.Log("price of", token1.Name, price)

}
