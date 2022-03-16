package uniswapv3

import (
	"context"
	"testing"
	"time"

	"github.com/netbound/dex-feed/token"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ethUsdcPool = common.HexToAddress("0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8")

	token1 = token.NewToken("WETH", wethAddress, 18)
	token0 = token.NewToken("USDC", usdcAddress, 6)
)

var p, _ = newTestPool()
var client *ethclient.Client

func newTestPool() (*Pool, error) {
	var err error

	i := PoolImmutables{
		Token0: token0,
		Token1: token1,
		Fee:    fee,
	}

	rpcApi := "http://localhost:8080/eth"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	client, err = ethclient.DialContext(ctx, rpcApi)
	if err != nil {
		return nil, err
	}

	p, err := NewPool(client, "WETHUSDC", ethUsdcPool, i)
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
	err := p.UpdateState(client)
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
	price, err := p.PriceOf(token1.Address)
	if err != nil {
		t.Fatalf("error getting price: %s", err)
	}

	t.Log("price of", token1.Name, price)
}

func BenchmarkPriceOf(b *testing.B) {
	_, err := p.PriceOf(token1.Address)
	if err != nil {
		b.Fatalf("error getting price: %s", err)
	}
}

func TestEncodeDecodePool(t *testing.T) {
	poolBytes, err := p.EncodePool()
	if err != nil {
		t.Fatal(err)
	}

	newPool, err := DecodePool(poolBytes)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%#v", newPool)
}
