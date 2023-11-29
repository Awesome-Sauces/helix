package account

import (
	"math/big"
)

func NumberToDrops(amount float64) big.Int {
	if amount < float64(0.000001) {
		return *big.NewInt(0)
	}

	return *big.NewInt(int64(amount * 1000000))
}
