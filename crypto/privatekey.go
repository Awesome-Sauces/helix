package crypto

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// PrivateKey represents a cryptographic private key.
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// SignToString signs a message and returns the base64-encoded signature as a string.
func (pk *PrivateKey) EncodeSignature(msg string) (string, error) {
	sig, err := pk.Sign(msg)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

// NewPrivateKey generates a new cryptographic private key.
func NewPrivateKey() (*PrivateKey, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: privKey}, nil
}

// NewPrivateKeyFromStr creates a new PrivateKey instance from a hexadecimal encoded string.
func DecodePrivateKey(encoded string) (*PrivateKey, error) {
	data, err := hex.DecodeString(strings.TrimPrefix(encoded, "0x"))
	if err != nil {
		return nil, err
	}
	privKey, err := crypto.ToECDSA(data)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: privKey}, nil
}

// Sign signs a message using the private key and returns the signature as a byte slice.
func (pk *PrivateKey) Sign(msg string) ([]byte, error) {
	hashedMsg := HashString(msg)
	signature, err := crypto.Sign(hashedMsg, pk.key)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// String returns the hexadecimal representation of the private key with "0x" prefix.
func (pk *PrivateKey) String() string {
	return "0x" + hex.EncodeToString(crypto.FromECDSA(pk.key))
}
