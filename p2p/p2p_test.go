package p2p

import (
	"encoding/binary"
	"fmt"
	"net"
	"testing"
)

func joinByteArrays(v1 []byte, v2 []byte) []byte {
	return append(v1, v2...)
}

func TestByteHeaders(t *testing.T) {
	bytes := []byte{}

	header := []byte("getClassMaxByteArrayLengthCapAll")

	bytes = append(bytes, []byte{0x00, 0x00, 0x00, 0x00}...)
	//bytes = append(bytes, 0xAF)

	bytes = joinByteArrays(bytes, header)

	for key, val := range bytes {
		fmt.Println(key, ":", val)
	}

	var sstring []byte

	for key, val := range bytes {
		if key < 2 {
			continue
		}

		sstring = append(sstring, val)
	}

	fmt.Println(string(sstring))
}

func TestClient(t *testing.T) {

	go StartTCPServer()

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return
	}
	defer conn.Close()

	// Sending 4 bytes
	message := []byte{0x00, 0x01, 0x02, 0x03}
	_, err = conn.Write(message)
	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		return
	}

	// Receiving 2 bytes back
	buffer := make([]byte, 2)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	fmt.Println("Received response:", buffer)

	/*
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
	*/

}

// sendLength sends the length of the message as a 4-byte integer in big-endian format
func sendLength(conn net.Conn, length int) error {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(length))
	_, err := conn.Write(buf)
	return err
}
