package p2p

import (
	"testing"
)

func TestClient(t *testing.T) {
	go Start("localhost:8080")
	go Start("localhost:8079")
	go Start("localhost:8077")

	Call("localhost:8080")
	Call("localhost:8079")
	Call("localhost:8077")

}
