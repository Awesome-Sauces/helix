// Package crypto provides utilities for working with cryptographic operations
// such as key generation, signing, and verification.
package crypto

import (
	"log"
)

// KeyPair represents a client with a private key, public key, and address.
type KeyPair struct {
	PrivateKey *PrivateKey
	PublicKey  *PublicKey
	Address    string
}

// Updated 0.0.1v

// NewKeyPair creates a new KeyPair with a generated private key.
func NewKeyPair() *KeyPair {
	privateKey, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	publicKey := privateKey.PublicKey()
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    publicKey.Address(),
	}
}

// KeyPairFromPK creates a new KeyPair from a private key string.
func KeyPairFromPK(key string) *KeyPair {
	privateKey, err := DecodePrivateKey(key)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	publicKey := privateKey.PublicKey()
	return &KeyPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Address:    publicKey.Address(),
	}
}

// Sign signs a transaction using the client's private key and returns the signature as a string.
func (skc *KeyPair) Sign(tx string) string {
	sig, err := skc.PrivateKey.EncodeSignature(tx)
	if err != nil {
		log.Fatal(err)
		return "nil"
	}
	return sig
}
