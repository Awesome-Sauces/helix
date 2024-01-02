package p2p

import (
	"fmt"
	"net"
	"os"
)

func Call(host string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		os.Exit(1)
	}

	message := []byte("Hello TCP server!")
	conn.Write(message)
	fmt.Println("Sent:", string(message))
	conn.Close()
}
