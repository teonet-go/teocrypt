// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
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
)

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
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
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
