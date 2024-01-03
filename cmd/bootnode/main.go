package main

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/Awesome-Sauces/helix/crypto/xdr"
)

func main() {
	// Hypothetical image header format.
	type ImageHeader struct {
		Signature   [3]byte
		Version     uint32
		IsGrayscale bool
		NumSections uint32
	}

	// XDR encoded data described by the above structure.  Typically this would
	// be read from a file or across the network, but use a manual byte array
	// here as an example.
	encodedData := []byte{
		0xAB, 0xCD, 0xEF, 0x00, // Signature
		0x00, 0x00, 0x00, 0x02, // Version
		0x00, 0x00, 0x00, 0x01, // IsGrayscale
		0x00, 0x00, 0x00, 0x0A, // NumSections
	}

	xdrString := base64.StdEncoding.EncodeToString(encodedData)

	fmt.Printf("string xdr: %s\n", xdrString)

	// Decode the base64-encoded XDR string into a byte slice
	decodedData, err := base64.StdEncoding.DecodeString(xdrString)
	if err != nil {
		// Handle the error
		fmt.Println(err)
	}

	if decodedData[7] == encodedData[7] {
		fmt.Println("IT WORKS")
	}

	// Declare a variable to provide Unmarshal with a concrete type and instance
	// to decode into.
	var h ImageHeader
	bytesRead, err := xdr.Unmarshal(bytes.NewReader(encodedData), &h)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("bytes read:", bytesRead)
	fmt.Printf("h: %+v", h)
}

/*
func main() {
	// Generate a mnemonic with a specific word count (e.g., 12 words)
	mnemonic, err := ethereum.GenerateMnemonic(12)
	if err != nil {
		log.Fatalf("Failed to generate mnemonic: %v", err)
	}
	fmt.Printf("Generated Mnemonic: %s\n", mnemonic)

	// Convert mnemonic to private key
	privateKey, err := ethereum.MnemonicToPrivateKey(mnemonic)
	if err != nil {
		log.Fatalf("Failed to create private key: %v", err)
	}

	// Ensure privateKey is not nil before using it
	if privateKey == nil {
		log.Fatal("Private key is nil")
	}

	// Get public key from the private key
	publicKey := privateKey.PublicKey()

	fmt.Println(privateKey.ToString())

	// Example: Convert public key to Ethereum address and print it
	address := ethereum.PublicKeyToAddress(publicKey)
	fmt.Printf("Ethereum Address: %s\n", address)

}
*/
