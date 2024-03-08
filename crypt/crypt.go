// Copyright 2023-24 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Crypt package contains functions to Encrypt and Decrypt data using
// Advanced Encryption Standard (AES).
package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

// ErrInvalidInputFile is returned when the input data is not valid.
var ErrInvalidInputFile = errors.New("invalid input file")

// Encrypt encrypts data using AES-GCM with a random nonce.
// It takes a key and data to encrypt as byte slices, and returns the
// ciphertext with nonce prepended and an error if one occurred.
func Encrypt(key, data []byte) (ciphertext []byte, err error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return
	}

	ciphertext = gcm.Seal(nonce, nonce, data, nil)

	return
}

// EncryptWriter creates stream cipher writer to encrypt output file.
func EncryptWriter(outputFile io.Writer, key []byte) (writer io.Writer, err error) {

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// Crete iv and write it to output file
	iv := make([]byte, block.BlockSize())
	if _, err = rand.Read(iv); err != nil {
		return
	}
	if _, err = outputFile.Write(iv); err != nil {
		return
	}

	// Create stream cipher
	stream := cipher.NewCTR(block, iv)

	// Create stream cipher writer
	writer = &cipher.StreamWriter{S: stream, W: outputFile}

	return
}

// Decrypt decrypts encrypted data using the provided key.
// It returns the decrypted plaintext and an error if one occurs.
func Decrypt(key, data []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, ErrInvalidInputFile
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// DecryptReader creates stream cipher reader to decrypt input file.
func DecryptReader(inputFile io.Reader, key []byte) (reader io.Reader, err error) {

	// Create cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// Read iv from input file
	iv := make([]byte, block.BlockSize())
	_, err = io.ReadFull(inputFile, iv)
	if err != nil {
		return
	}

	// Create stream cipher
	stream := cipher.NewCTR(block, iv)

	// Create stream cipher reader
	reader = &cipher.StreamReader{S: stream, R: inputFile}

	return
}

// EncryptXor encrypts the given data using XOR with the given key.
// The key is repeated to match the length of the data.
// Returns the encrypted ciphertext.
func EncryptXor(key, data []byte) []byte {

	ciphertext := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		ciphertext[i] = data[i] ^ key[i%len(key)]
	}

	return ciphertext
}

// DecryptXor decrypts the given data using XOR with the given key.
// The key is repeated to match the length of the data.
// Returns the decrypted plaintext.
func DecryptXor(key, data []byte) []byte {

	plaintext := make([]byte, len(data))

	for i := 0; i < len(data); i++ {
		plaintext[i] = data[i] ^ key[i%len(key)]
	}

	return plaintext
}

// GenerateKey generates a cryptographically secure random 32 byte key.
// It uses the rand package to generate the random bytes and returns
// the 32 byte key slice. It returns an error if there is a problem
// reading from the system's random number generator.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// HashKey hashes the given password string using SHA256 and returns the hash
// as a byte slice.
func HashKey(passwd string) []byte {
	h := sha256.New()
	h.Write([]byte(passwd))
	return h.Sum(nil)
}
