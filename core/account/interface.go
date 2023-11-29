package account

import "math/big"

type Accounter interface {
	GetBalance() big.Int

	Credit(amount big.Int) error
	Debit(amount big.Int) error
	SetBalance(amount big.Int) error

	GetAddress() string
}
