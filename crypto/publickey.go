package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// PublicKey represents a cryptographic public key.
type PublicKey struct {
	key *ecdsa.PublicKey
}

// String returns the hexadecimal representation of the public key with "0x" prefix.
func (pub *PublicKey) String() string {
	return "0x" + hex.EncodeToString(crypto.CompressPubkey(pub.key))
}

// PublicKey returns the corresponding public key.
func (pk *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{key: &pk.key.PublicKey}
}

// Address returns the Ethereum-style address derived from the public key.
func (pub *PublicKey) Address() string {
	addressBytes := crypto.PubkeyToAddress(*pub.key).Bytes()
	return "0x" + hex.EncodeToString(addressBytes)
}

// Verify verifies a message signature using the public key.
func (pub *PublicKey) VerifySignature(msg string, signature []byte) bool {
	hashedMsg := HashString(msg)

	// Convert the public key to a byte slice
	pubKeyBytes := elliptic.Marshal(crypto.S256(), pub.key.X, pub.key.Y)

	valid := crypto.VerifySignature(pubKeyBytes, hashedMsg, signature)
	return valid
}

// NewPublicKeyFromStr creates a new PublicKey instance from a hexadecimal encoded string.
func DecodePublicKey(encoded string) (*PublicKey, error) {
	data, err := hex.DecodeString(strings.TrimPrefix(encoded, "0x"))
	if err != nil {
		return nil, err
	}
	var pubKey ecdsa.PublicKey
	pubKey.X, pubKey.Y = new(big.Int).SetBytes(data[:32]), new(big.Int).SetBytes(data[32:])
	return &PublicKey{key: &pubKey}, nil
}

// RecoverPublicKeyFromSignature recovers a public key from a message signature.
func RecoverPublicKeyFromSignature(message string, signature string) (*PublicKey, error) {
	sigBytes, err := base64.StdEncoding.DecodeString(signature) // Decode base64 signature
	if err != nil {
		return nil, err
	}
	hashedMsg := HashString(message)
	pubKey, err := crypto.SigToPub(hashedMsg, sigBytes)
	if err != nil {
		return nil, err
	}
	return &PublicKey{key: pubKey}, nil
}
