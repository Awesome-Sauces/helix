package tcpp

import (
	"fmt"
	"log"
	"net"
)

// first byte = request size in bytes
// second byte = request size in bytes
// third byte = request size in bytes
// fourth byte = request size in bytes
// result first_byte + second_byte + third_byte + fourth_byte = max size request = 1020 bytes

// fifth byte-thirty fifth byte = endpoint name (Similar to HTTP)
// endpoint has a max length of 32 char(s) = 32 bytes
// The rest of the request is from the 36 byte and on
// the endpoint will handle the processing with the help of
// easy

// We assume anything beyond byte 35 are typed values
// With 6 bytes for the type, the types are the following
// str(Until end byte is defined as 0x00 (Null)), u32 & i32 (4 bytes),
// u64 & i64 (8 bytes), u128 & i128 (16 bytes),
// variable names are given 16 bytes and come before type declaration.
// variable names and type declaration are seperated by the ascii null (0x00)
func DigestRequest(conn net.Conn) (string, []byte) {
	defer conn.Close()

	buffer := make([]byte, 4) // Buffer to read 4 bytes
	_, err := conn.Read(buffer)
	if err != nil {
		log.Println("Error reading:", err.Error())
		return "nil", nil
	}

	size := buffer[0] + buffer[1] + buffer[2] + buffer[3]

	buffer = make([]byte, size)
	_, err = conn.Read(buffer)
	if err != nil {
		log.Println("Error reading:", err.Error())
		return "nil", nil
	}

	b_endpoint := []byte{}

	t_buffer := []byte{}

	for i, val := range buffer {
		if i <= 35 && i > 3 {
			if val != 0 {
				b_endpoint = append(b_endpoint, val)
			}
		}

		t_buffer = append(t_buffer, val)
	}

	return string(b_endpoint), t_buffer
}

func SearchValues() {

	// An example of transmitting variables over tcpp (TCP plus)
	values := []byte{
		0x00,
		0x74,
		0x6F,
		0x6B,
		0x65,
		0x6E,
		0x00,
		0x69,
		0x36,
		0x34,
		0x00,
		0x2E,
		0x49,
		0xDC,
		0x3F,
		0x0B,
		0xEF,
		0xDD,
		0xEE, // Since the 8 bytes for i64 are over here we can just start
		0x6E, // Another variable since we auto recognize the start of a new one
		0x61,
		0x6D,
		0x65,
		0x00,
		0x73,
		0x74,
		0x72,
		0x00, // New variable call `name` with type `str` or (string)
		0x78,
		0x72,
		0x70,
		0x00, // After we end this string we will add an extra byte
		0x04, // This is ASCII *EOT* which means the variables end
		// It also means End of Transmission.
	}

	position := 0

	forward := func() {
		if position+1 <= len(values) {
			position++
		}
	}

	fmt.Println(forward)
	fmt.Println(position)

	for i := 0; i < len(values); i++ {
		
	}

}
