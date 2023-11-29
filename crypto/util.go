package crypto

import (
	"crypto/sha256"
	"encoding/base64"
)

// StringToSignature decodes a hexadecimal encoded signature string into a byte slice.
func DecodeSignature(encodedSig string) ([]byte, error) {
	sigBytes, err := base64.StdEncoding.DecodeString(encodedSig)
	return sigBytes, err
}

// hashString computes the SHA-256 hash of the input string and returns the result as a byte slice.
func HashString(data string) []byte {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashedMessage := hasher.Sum(nil)
	return hashedMessage
}
