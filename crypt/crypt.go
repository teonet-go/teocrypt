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

var ErrInvalidInputFile = errors.New("invalid input file")

// Encrypt encrypts data by key.
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

// Decrypt decrypts data by key.
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

// GenerateKey generates random key.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// HashKey creates a new key from password.
func HashKey(passwd string) []byte {
	h := sha256.New()
	h.Write([]byte(passwd))
	return h.Sum(nil)
}
