// teocrypt.go
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
)

const (
	appShort   = "teocrypt"
	appName    = "Teonet encrypt/decrypt text application"
	appLong    = ""
	appVersion = "0.0.1"
)

func main() {

	// Application logo
	fmt.Println(appName + " ver " + appVersion)

	// Parse application command line parameters
	var save, decrypt bool
	var text, passwd string
	flag.StringVar(&text, "t", "", "text to encrypt/decrypt")
	flag.StringVar(&passwd, "p", "", "password used to encrypt/decrypt")
	flag.BoolVar(&decrypt, "d", decrypt, "decrypt text specified in -t flag")
	flag.BoolVar(&save, "save-password", save, "save password specified in -p flag on this device")
	flag.Parse()

	// Get key
	var err error
	var key []byte
	if len(passwd) > 0 {
		key = HashKey(passwd)
	} else if key, err = GenerateKey(); err != nil {
		fmt.Printf("can't generate new key, error: %s\n", err)
		return
	}

	// Get data
	var data []byte
	if len(text) > 0 {
		data = []byte(text)
	} else {
		data = []byte("super secret text")
	}

	// Encrypt input data by key
	if !decrypt {
		ciphertext, err := Encrypt(key, data)
		if err != nil {
			fmt.Printf("can't encode input text, error: %s\n", err)
			return
		}
		fmt.Printf("encrypted text:\n%s\n", hex.EncodeToString(ciphertext))
		return
	}

	// Decrypt input data by key
	data, err = hex.DecodeString(string(data))
	if err != nil {
		fmt.Printf("can't decrypt input text, error: %s\n", err)
		return
	}
	plaintext, err := Decrypt(key, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("plain text:\n%s\n", plaintext)
}

// Encrypt encrypts data by key
func Encrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt decrypts data by key
func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateKey generates random key
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// HashKey creates a new key from password
func HashKey(passwd string) []byte {
	h := sha256.New()
	h.Write([]byte(passwd))
	return h.Sum(nil)
}
