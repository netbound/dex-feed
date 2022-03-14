package uniswapv3

import (
	"context"
	"errors"
	"math/big"
	"time"

	univ3pool "dex-feed/bindings/uniswap_v3/pool"
	"dex-feed/token"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrWrongToken = errors.New("uniswap v3 pool: PriceOf: token not in pool")
)

type PoolImmutables struct {
	Token0 token.Token
	Token1 token.Token
	Fee    int64
}

type PoolState struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
}

type Pool struct {
	Name    string
	Address common.Address

	Caller *univ3pool.Univ3poolCaller

	Immutables PoolImmutables
	State      PoolState // The last known Pool State
}

func NewPool(client *ethclient.Client, name string, poolAddress common.Address, immutables PoolImmutables) (*Pool, error) {
	caller, err := univ3pool.NewUniv3poolCaller(poolAddress, client)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Name:       name,
		Address:    poolAddress,
		Caller:     caller,
		Immutables: immutables,
		// Initialize empty PoolState
		State: PoolState{},
	}, nil
}

// UpdateState updates the internal pool state. Should be called every time the state changes on-chain
// i.e. on a new block.
func (p *Pool) UpdateState() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	opts := &bind.CallOpts{Context: ctx}
	slot0, err := p.Caller.Slot0(opts)
	if err != nil {
		return err
	}

	p.State.SqrtPriceX96 = slot0.SqrtPriceX96
	p.State.Tick = slot0.Tick

	return nil
}

func (p *Pool) PriceOf(token token.Token) (float64, error) {
	var (
		token0Multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(p.Immutables.Token0.Decimals), nil)
		token1Multiplier = new(big.Int).Exp(big.NewInt(10), big.NewInt(p.Immutables.Token1.Decimals), nil)
	)

	if token.Address == p.Immutables.Token0.Address {
		numerator := new(big.Int).Exp(p.State.SqrtPriceX96, big.NewInt(2), nil)
		// multiply by token decimals
		numerator = numerator.Mul(numerator, token0Multiplier)
		n := new(big.Float).SetInt(numerator)

		denominator := new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)
		d := new(big.Float).SetInt(denominator)

		res := n.Quo(n, d)
		price, _ := res.Quo(res, new(big.Float).SetInt(token1Multiplier)).Float64()
		return price, nil
	} else if token.Address == p.Immutables.Token1.Address {
		numerator := new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil)
		numerator = numerator.Mul(numerator, token1Multiplier)
		n := new(big.Float).SetInt(numerator)

		denominator := new(big.Int).Exp(p.State.SqrtPriceX96, big.NewInt(2), nil)
		d := new(big.Float).SetInt(denominator)

		res := n.Quo(n, d)
		price, _ := res.Quo(res, new(big.Float).SetInt(token0Multiplier)).Float64()

		return price, nil
	}

	return 0, ErrWrongToken
}
