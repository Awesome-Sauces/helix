package hydrogen

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
)

func TestMapVariables(t *testing.T) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(1234567891234567890))
	values, err := BareEncode([]byte("token"), []byte("i64"), b, []byte("name"), []byte("str"), []byte("xrp"), []byte("address"), []byte("str"), []byte("0xA1FC67E"), []byte("send_amount"), []byte("i64"), b)

	if err != nil {
		log.Fatal(err)
	}
	pass_back := MapVariables(values)

	for key, val := range pass_back {
		log.Println(key, ":", val.Name, ":", val.Type, ":", val.Bytes)
	}
}

func TestServerCore(t *testing.T) {
	StartTCPServer()

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(1234567891234567890))
	values, err := BareEncode([]byte("token"), []byte("i64"), b, []byte("name"), []byte("str"), []byte("xrp"), []byte("address"), []byte("str"), []byte("0xA1FC67E"), []byte("send_amount"), []byte("i64"), b)

	if err != nil {
		log.Fatal(err)
	}

	request, err := SendRequest("StartPropagation", values)

	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	conn.Write(request)
	conn.Close()
}
