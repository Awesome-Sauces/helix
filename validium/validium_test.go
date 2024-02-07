package validium

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/Awesome-Sauces/helix/crypto/ethereum"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/crypto"
)

type Sig struct {
	Signature []byte
	Address   string
	ID        int
}

type Transaction struct {
	Sender     string
	Payload    []byte
	Originator int
	Signatures []Sig
}

type Node struct {
	ID       int
	Listener chan Transaction
}

func StructSize(s interface{}) int {
	size := 0
	value := reflect.ValueOf(s)
	//typeOf := reflect.TypeOf(s)

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		//fieldType := typeOf.Field(i)

		switch field.Kind() {
		case reflect.Slice:
			sliceElemType := field.Type().Elem()
			sliceElemSize := int(unsafe.Sizeof(sliceElemType.Align()))
			sliceLen := field.Len()
			sliceElemSize *= sliceLen
			size += sliceElemSize

		default:
			size += int(unsafe.Sizeof(field.Addr().Interface()))
		}
	}

	return size
}

func (n *Node) Listen(handler func(bytes Transaction)) {
	for {
		conn := <-n.Listener

		handler(conn)
	}
}

func (n Node) Dial(node *Node, input Transaction) {
	node.Listener <- input
}

func TestConsensus(t *testing.T) {
	node1 := &Node{
		ID:       5,
		Listener: make(chan Transaction),
	}

	node2 := &Node{
		ID:       3,
		Listener: make(chan Transaction),
	}

	handler := func(bytes Transaction) {
		fmt.Println(bytes.Originator)
		fmt.Println(bytes.Sender)
		fmt.Println(bytes.Payload)
		fmt.Println(bytes.Signatures[0].Address)
		fmt.Println(bytes.Signatures[0].ID)
		fmt.Println(bytes.Signatures[0].Signature)
	}

	privKey, _ := ethereum.NewPrivateKey()
	pubKey := privKey.PublicKey()
	address := ethereum.PublicKeyToAddress(pubKey)

	go node1.Listen(handler)

	data := []byte(fmt.Sprintf("%s SENDS 50 TO %s", address, address))
	hash := crypto.Keccak256Hash(data)
	signature, err := privKey.Sign(hash.Bytes())

	if err != nil {
		fmt.Println(err)
		fmt.Println("HOLD UP")
		return
	}

	node2.Dial(node1, Transaction{
		Sender:     address,
		Payload:    data,
		Originator: 13,
		Signatures: []Sig{{Signature: signature, Address: address, ID: 13}},
	})

	tx := Transaction{
		Sender:     address,
		Payload:    data,
		Originator: 13,
		Signatures: []Sig{{Signature: signature, Address: address, ID: 13}},
	}

	fmt.Println("Size of Transaction: ", StructSize(tx))

	fmt.Println("SIZER: ", unsafe.Sizeof(tx))

}
