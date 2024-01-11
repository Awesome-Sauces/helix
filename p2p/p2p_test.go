package p2p

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"testing"
)

func joinByteArrays(v1 []byte, v2 []byte) []byte {
	return append(v1, v2...)
}

func CompileToTCPP(args ...[]byte) ([]byte, error) {
	if len(args)%3 != 0 {
		return nil, errors.New("INVALID: Amount of byte arrays provided insufficient")
	}

	request := []byte{}

	for i := 0; i < len(args)/3; i++ {

		varn := args[(1+(i*((1%(i+1))*3)))-1]
		vart := args[(2+(i*((1%(i+1))*3)))-1]
		varv := args[(3+(i*((1%(i+1))*3)))-1]

		request = append(request, varn...)
		request = append(request, 0x00)
		request = append(request, vart...)
		request = append(request, 0x00)
		request = append(request, varv...)

		if string(vart) == "str" {
			request = append(request, 0x00)
		}

	}

	request = append(request, 0x04)

	return request, nil
}

func TestVariableSearch(t *testing.T) {
	// An example of transmitting variables over tcpp (TCP plus)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(1234567891234567890))
	values, err := CompileToTCPP([]byte("token"), []byte("i64"), b, []byte("name"), []byte("str"), []byte("xrp"))

	if err != nil {
		fmt.Println(err)
		return
	}

	i := 0

	forward := func() {
		if i+1 <= len(values) {
			i++
		}
	}

	fmt.Println(i)

	if i >= len(values) {
		return
	}

	build_var := func() ([]byte, []byte, []byte) {
		var_name := []byte{}
		var_type := []byte{}
		var_value := []byte{}

		for values[i] != 0x00 {
			var_name = append(var_name, values[i])
			forward()
		}

		forward()

		for values[i] != 0x00 {
			var_type = append(var_type, values[i])
			forward()
		}

		forward()

		if string(var_type) == "i32" || string(var_type) == "u32" {
			for iter := 0; iter < 4; iter++ {
				var_value = append(var_value, values[i])
				forward()
			}
		}

		if string(var_type) == "i64" || string(var_type) == "u64" {
			for iter := 0; iter < 8; iter++ {
				var_value = append(var_value, values[i])
				forward()
			}
		}

		if string(var_type) == "i128" || string(var_type) == "u128" {
			for iter := 0; iter < 16; iter++ {
				var_value = append(var_value, values[i])
				forward()
			}
		}

		if string(var_type) == "str" {
			for values[i] != 0x00 {
				var_value = append(var_value, values[i])
				forward()
			}
		}

		return var_name, var_type, var_value
	}

	vn, vt, vv := build_var()

	fmt.Println(string(vn))
	fmt.Println(string(vt))
	fmt.Println(int64(binary.LittleEndian.Uint64(vv)))

	varn, vart, varv := build_var()

	fmt.Println(string(varn))
	fmt.Println(string(vart))
	fmt.Println(string(varv))

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
