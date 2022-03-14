package token

import "github.com/ethereum/go-ethereum/common"

type Token struct {
	Name     string
	Address  common.Address
	Decimals int64
}

func NewToken(name string, address common.Address, decimals int64) Token {
	return Token{
		Name:     name,
		Address:  address,
		Decimals: decimals,
	}
}
