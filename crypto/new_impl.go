package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/tyler-smith/go-bip39" // BIP39 library
)

// PrivateKey struct represents a private key
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// NewPrivateKey generates a new ECDSA private key
func NewPrivateKey() (*PrivateKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: privKey}, nil
}

// Sign signs data using the private key
func (pk *PrivateKey) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	signature, err := ecdsa.SignASN1(rand.Reader, pk.key, hash[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// GenerateMnemonic generates a new mnemonic phrase
func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(160)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func MnemonicToPrivateKey(mnemonic string) (*PrivateKey, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, "") // You can use a passphrase here

	// Convert seed (byte slice) to an io.Reader
	seedReader := bytes.NewReader(seed)

	// Now use seedReader as the io.Reader for ecdsa.GenerateKey
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), seedReader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: privKey}, nil
}

// String returns the hexadecimal representation of the private key
func (pk *PrivateKey) ToString() string {
	return hex.EncodeToString(pk.key.D.Bytes())
}
