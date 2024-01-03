package p2p

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

// Config struct to hold server configuration
type Config struct {
	PrivateKey string `json:"private_key"` // Ethereum-compatible private key
	// "IP>PUBLIC_KEY"
	Peers       []string `json:"peers"`        // List of peers
	LatestBlock string   `json:"latest_block"` // Identifier for the latest block
}

// Encrypts the configuration data using AES
func encryptConfigData(config Config, passphrase string) ([]byte, error) {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	key := sha256.Sum256([]byte(passphrase))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(jsonData))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], jsonData)

	return ciphertext, nil
}

// Decrypts the configuration data using AES
func decryptConfigData(ciphertext []byte, passphrase string) (Config, error) {
	var config Config

	key := sha256.Sum256([]byte(passphrase))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return config, err
	}

	if len(ciphertext) < aes.BlockSize {
		return config, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	err = json.Unmarshal(ciphertext, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {
	configFile := "path/to/config.file"
	passphrase := getPassphraseFromUser() // Implement this function to securely obtain the passphrase

	// Read the encrypted config file
	encryptedConfig, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	// Decrypt the config file
	config, err := decryptConfigData(encryptedConfig, passphrase)
	if err != nil {
		fmt.Printf("Error decrypting config file: %v\n", err)
		os.Exit(1)
	}

	// Initialize and start your TCP server using the configuration
	// ...
	fmt.Println("Server configuration loaded:", config)
	// For example, start listening on a TCP port
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Error setting up TCP listener: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("TCP server listening on localhost:8080")

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn, &config)
	}
}

func handleConnection(conn net.Conn, config *Config) {
	// Handle the connection
	// For example, read data from the connection and process it
	// ...

	conn.Close()
}

func getPassphraseFromUser() string {
	// Implement a secure way to get the passphrase from the user
	// This is a placeholder implementation
	return "your_secure_passphrase"
}
