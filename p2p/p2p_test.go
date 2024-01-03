package p2p

import (
	"log"
	"testing"

	"github.com/Awesome-Sauces/helix/crypto/ethereum"
)

func TestClient(t *testing.T) {
	go Start("localhost:8080")
	go Start("localhost:8079")
	go Start("localhost:8077")

	Call("localhost:8080")
	Call("localhost:8079")
	Call("localhost:8077")

	pk, err := ethereum.NewPrivateKey()

	if err != nil {
		return
	}

	log.Println(ethereum.PublicKeyToAddress(pk.PublicKey()))

}
