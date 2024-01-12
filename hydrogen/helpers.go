package hydrogen

import (
	"errors"
	"fmt"
)

func SendRequest(endpoint string, args ...[]byte) ([]byte, error) {
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

	req_size := len(request) + len([]byte(endpoint)) + 4

	if req_size > 255*4 {
		return []byte{}, errors.New(fmt.Sprintf("INVALID: Request is bigger than allowed (%d > 1020)", req_size))
	}

	// Calculate the 4 bytes
	size := []byte{
		byte(req_size & 0xFF),
		byte((req_size >> 8) & 0xFF),
		byte((req_size >> 16) & 0xFF),
		byte((req_size >> 24) & 0xFF),
	}

	h_request := []byte{}

	h_request = append(h_request, size...)
	h_request = append(h_request, []byte(endpoint)...)
	h_request = append(h_request, request...)

	return h_request, nil
}

func BareEncode(args ...[]byte) ([]byte, error) {
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
