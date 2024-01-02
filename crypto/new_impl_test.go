package crypto

import (
	"fmt"
	"log"
	"testing"
)

func TestCrypt(t *testing.T) {
	pk, err := NewPrivateKey()

	if err != nil {
		log.Println(err)
	}

	pk.ToString()

	mnemonic, err := GenerateMnemonic()

	fmt.Println(mnemonic)
	if err != nil {
		log.Println(err)
	}
}
