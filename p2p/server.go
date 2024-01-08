package p2p

import (
	"bytes"
	"fmt"
	"net"

	"github.com/Awesome-Sauces/helix/crypto/xdr"
)

func StartTCPServer() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		go handleTraffic(conn)
	}
}

func handleTraffic(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4) // Buffer to read 4 bytes
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}

	h := "Send"

	var w bytes.Buffer
	bytesWritten, err := xdr.Marshal(&w, &h)
	if err != nil {
		fmt.Println(err)
		return
	}

	encodedData := w.Bytes()
	fmt.Println("bytes written:", bytesWritten)
	fmt.Println("encoded data:", encodedData)

	// Process the data (optional, based on your logic)
	fmt.Println("Received data:", buffer)

	// Sending back 2 bytes
	response := []byte{0x01, 0x02}
	conn.Write(response)
}
