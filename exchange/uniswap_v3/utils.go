package uniswapv3

import (
	"math"

	"github.com/ethereum/go-ethereum/common"
)

// Gets the amounts out for the current price and lpValue, given the lower and upper price bounds
// of the position.
func GetAmountsOut(price, lpValue, lowerPrice, upperPrice float64) (float64, float64) {
	l := 1 / (math.Sqrt(price) - math.Sqrt(lowerPrice))

	x := l / ((math.Sqrt(upperPrice) * math.Sqrt(price)) / (math.Sqrt(upperPrice) - math.Sqrt(price)))
	y := lpValue / (1 + x*price)
	x = (lpValue - y) / price

	return x, y
}

// Gets the initialized tick at a price given the tickspace.
// 5bps -> 10, 30bps -> 60, 100bps -> 200
func GetInitializedTickAtPrice(price float64, tickSpace int) int64 {
	raw := int(math.Floor(math.Log(price) / math.Log(1.0001)))

	if raw%tickSpace < (tickSpace / 2) {
		return int64(raw - raw%tickSpace)
	}

	return int64(raw + (60 - raw%tickSpace))
}

// Gets the tick at the given price.
func GetTickAtPrice(price float64) int64 {
	return int64(math.Floor(math.Log(price) / math.Log(1.0001)))
}

// Gets the price at the given tick.
func GetPriceAtTick(tick int64) float64 {
	return math.Pow(1.0001, float64(tick))
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
