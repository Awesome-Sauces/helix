package account

import (
	"log"
	"testing"
)

func TestStruct(t *testing.T) {
	account := BlankAccount("0xA1SAUCE")
	//.000001
	NumberToDrops(1)

	account.Credit(NumberToDrops(1))

	log.Println(account.GetBalance().String())
}
