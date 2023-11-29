package account

import (
	"errors"
	"math/big"
)

// Account struct used accross the whole code
// Balance is defined in units of "drops"

// Similar approach to XRPL
// A unit is defined by drops (10 drops = 0.000001 of any unit)
// One whole unit would be 1,000,000 drops
type Account struct {
	address string
	balance *big.Int
}

// Interfacing
func (account Account) GetBalance() *big.Int {
	return account.balance
}

func (account Account) GetAddress() string {
	return account.address
}

// Gives account amount of balance
func (account *Account) Credit(amount big.Int) error {

	account.balance = big.NewInt(0).Add(account.balance, &amount)

	return nil
}

// Takes amount of balance from account
func (account *Account) Debit(amount big.Int) error {
	balance := big.NewInt(0).Sub(account.balance, &amount)

	if balance.Cmp(big.NewInt(0)) == -1 {
		return errors.New("INSUFFICIENT BALANCE")
	}

	account.SetBalance(balance)

	return nil
}

// For rudimentary operations
func (account *Account) SetBalance(amount *big.Int) {
	account.balance = amount
}

// Creates an Account struct that holds no balance
func BlankAccount(address string) *Account {
	return &Account{
		address: address,
		balance: big.NewInt(0),
	}
}
