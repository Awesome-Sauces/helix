package hydrogen

import (
	"errors"
	"log"
	"net"
)

const (
	NULL        = 0
	String      = 1
	Signed32    = 2
	Unsigned32  = 3
	Signed64    = 4
	Unsigned64  = 5
	Signed128   = 6
	Unsigned128 = 7
)

type FunctionArgument struct {
	Name  string
	Type  string
	Bytes []byte
}

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
		return "FAILED TO READ BYTES (MIGHT HAVE BYTE OVERFLOW/UNDERFLOW)", nil
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

func MapVariables(data []byte) []FunctionArgument {
	// Example use case
	//b := make([]byte, 8)
	//binary.LittleEndian.PutUint64(b, uint64(1234567891234567890))
	//values, err := BareEncode([]byte("token"), []byte("i64"), b, []byte("name"), []byte("str"), []byte("xrp"))

	i := 0

	forward := func() bool {
		if i+1 <= len(data) {
			i++
			return false
		}

		return true
	}

	if i >= len(data) {
		return nil
	}

	build_var := func() ([]byte, []byte, []byte, error) {
		var_name := []byte{}
		var_type := []byte{}
		var_value := []byte{}

		for i != len(data) && data[i] != 0x00 {
			var_name = append(var_name, data[i])
			if forward() {
				return nil, nil, nil, errors.New("Variable Search ENDED")
			}
		}

		if forward() {
			return nil, nil, nil, errors.New("Variable Search ENDED")
		}

		for i != len(data) && data[i] != 0x00 {
			var_type = append(var_type, data[i])
			if forward() {
				return nil, nil, nil, errors.New("Variable Search ENDED")
			}
		}

		if forward() {
			return nil, nil, nil, errors.New("Variable Search ENDED")
		}

		if string(var_type) == "i32" || string(var_type) == "u32" {
			for iter := 0; iter < 4; iter++ {
				var_value = append(var_value, data[i])
				if forward() {
					return nil, nil, nil, errors.New("Variable Search ENDED")
				}
			}
		}

		if string(var_type) == "i64" || string(var_type) == "u64" {
			for iter := 0; iter < 8; iter++ {
				var_value = append(var_value, data[i])
				if forward() {
					return nil, nil, nil, errors.New("Variable Search ENDED")
				}
			}
		}

		if string(var_type) == "i128" || string(var_type) == "u128" {
			for iter := 0; iter < 16; iter++ {
				var_value = append(var_value, data[i])
				if forward() {
					return nil, nil, nil, errors.New("Variable Search ENDED")
				}
			}
		}

		if string(var_type) == "str" {
			for data[i] != 0x00 {
				var_value = append(var_value, data[i])
				if forward() {
					return nil, nil, nil, errors.New("Variable Search ENDED")
				}
			}
		}

		return var_name, var_type, var_value, nil
	}

	retval := []FunctionArgument{}

	for {
		vn, vt, vv, err := build_var()

		if err != nil {
			return retval
		}

		retval = append(retval, FunctionArgument{Name: string(vn), Type: string(vt), Bytes: vv})
	}

	//fmt.Println(string(vn))
	//fmt.Println(string(vt))
	//fmt.Println(int64(binary.LittleEndian.Uint64(vv)))
}
