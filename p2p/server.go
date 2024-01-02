package p2p

import (
	"fmt"
	"net"
	"os"
)

func Start(host string) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on ", host)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buffer := make([]byte, 1024)
	length, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received:", string(buffer[:length]))
	conn.Close()
}
