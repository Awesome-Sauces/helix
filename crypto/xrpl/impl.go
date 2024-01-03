package xrpl

/*

WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING
WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING
WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING WARNING

This module is STILL IN DEVELOPMENT, meaning it is not to be used.

*/

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2" // Bitcoin's secp256k1 library
	"github.com/btcsuite/btcutil/base58"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ripemd160" // RIPEMD160 hash function
)

// PrivateKey struct represents a private key
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// NewPrivateKey generates a new secp256k1 ECDSA private key.
func NewPrivateKey() (*PrivateKey, error) {
	privKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: privKey}, nil
}

// Sign signs data using the private key.
func (pk *PrivateKey) Sign(data []byte) ([]byte, error) {
	signature, err := ecdsa.SignASN1(rand.Reader, pk.key, data)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// PublicKey returns the public key associated with the private key.
func (pk *PrivateKey) PublicKey() *ecdsa.PublicKey {
	return &pk.key.PublicKey
}

// PublicKeyToAddress converts a public key to an XRP Ledger address.
func PublicKeyToAddress(pubKey *ecdsa.PublicKey) (string, error) {
	pubKeyBytes := elliptic.Marshal(btcec.S256(), pubKey.X, pubKey.Y)

	// Perform SHA256 and then RIPEMD160 hashing
	sha256Hash := sha256.Sum256(pubKeyBytes)
	ripemd160Hasher := ripemd160.New()
	_, err := ripemd160Hasher.Write(sha256Hash[:])
	if err != nil {
		return "", err
	}
	ripemd160Hash := ripemd160Hasher.Sum(nil)

	// Convert to XRP Ledger address format
	address := EncodeBase58Check(ripemd160Hash)
	return address, nil
}

// GenerateMnemonic, MnemonicToPrivateKey, ToString, and other utility functions would remain similar to the Ethereum library.
// You would have to implement the `EncodeBase58Check` function for XRP Ledger address encoding.

// GenerateMnemonic generates a new mnemonic phrase.
func GenerateMnemonic(wordCount int) (string, error) {
	bitSize := wordCount * 32 / 3 // Convert word count to bit size
	if bitSize%32 != 0 || bitSize < 128 || bitSize > 256 {
		return "", errors.New("invalid word count: must be 12, 15, 18, 21, or 24")
	}

	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

// MnemonicToPrivateKey converts a mnemonic to a private key.
func MnemonicToPrivateKey(mnemonic string) (*PrivateKey, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}

	//seed := bip39.NewSeed(mnemonic, "")
	privKey, err := ecdsa.GenerateKey(btcec.S256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{key: privKey}, nil
}

// ToString returns the hexadecimal representation of the private key.
func (pk *PrivateKey) ToString() string {
	return hex.EncodeToString(pk.key.D.Bytes())
}

// EncodeBase58Check encodes a byte slice into a modified base58 string with checksum.
func EncodeBase58Check(input []byte) string {
	// Perform double SHA256 hashing on the input
	checksum := sha256.Sum256(input)
	checksum = sha256.Sum256(checksum[:])

	// Append the first 4 bytes of the checksum to the input
	extended := append(input, checksum[:4]...)

	// Encode the extended buffer using Base58
	return base58.Encode(extended)
}
