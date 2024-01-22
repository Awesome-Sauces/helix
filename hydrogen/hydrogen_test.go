package hydrogen

import (
	"encoding/binary"
	"log"
	"testing"
)

func TestMapVariables(t *testing.T) {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(1234567891234567890))
	values, err := BareEncode([]byte("token"), []byte("i64"), b, []byte("name"), []byte("str"), []byte("xrp"), []byte("address"), []byte("str"), []byte("0xA1FC67E"))

	if err != nil {
		log.Fatal(err)
	}
	pass_back := MapVariables(values)

	for key, val := range pass_back {
		log.Println(key, ":", val.Name, ":", val.Type, ":", val.Bytes)
	}
}
