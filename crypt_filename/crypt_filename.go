// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Crypt filename package contains functions to Encrypt and Decrypt S3
// compatible file path included filenames.
package crypt_filename

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"
	"strings"

	"github.com/teonet-go/teocrypt/crypt"
)

// CryptFilename contains methods to encrypt and decrypt S3 filenames.
type CryptFilename struct {
	zipping      bool   // zip file names
	hashKey      []byte // hash key
	encryptFirst bool   // encrypt first folder
}

// New creates new CryptFilename object.
//
// Arguments description:
//
//	key - string key or password used to encrypt and decrypt filenames
//	zip - archive long path folders and filename if true
//	encryptFirst - encrypt first folder in path if true (false if omitted)
func New(key string, zip bool, encryptFirst ...bool) *CryptFilename {
	return &CryptFilename{
		hashKey:      crypt.HashKey(key),
		zipping:      zip,
		encryptFirst: len(encryptFirst) > 0 && encryptFirst[0],
	}
}

func (c CryptFilename) base64EncodeEscape(data []byte) string {
	r := strings.NewReplacer("/", "_")
	return r.Replace(base64.StdEncoding.EncodeToString(data))
}

func (c CryptFilename) base64DecodeEscape(str string) ([]byte, error) {
	r := strings.NewReplacer("_", "/")
	return base64.StdEncoding.DecodeString(r.Replace(str))
}

func (c CryptFilename) zip(data []byte) []byte {
	if !c.zipping {
		return data
	}

	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(data)
	w.Close()

	if len(data) <= b.Len() {
		return data
	}
	return b.Bytes()
}

func (c CryptFilename) unzip(data []byte) ([]byte, error) {
	if !c.zipping {
		return data, nil
	}

	var b bytes.Buffer
	t := bytes.NewBuffer(data)
	r, err := gzip.NewReader(t)
	if err != nil {
		return data, err
	}
	_, err = io.Copy(&b, r)
	if err != nil {
		return data, err
	}
	r.Close()

	return b.Bytes(), nil
}

// Encrypt encrypts s3 compatible full filename.
func (c CryptFilename) Encrypt(s string) (res string, err error) {
	parts := strings.Split(s, "/")
	for i, p := range parts {
		// Don't encrypt first folder of path
		if !c.encryptFirst && i == 0 {
			res = p
			continue
		}

		// Path slash
		if i > 0 {
			res += "/"
		}

		// Don't encrypt empy path
		if len(p) == 0 {
			continue
		}

		// Encrypt
		data := c.zip([]byte(p))
		ciphertext, err := crypt.Encrypt(c.hashKey, data)
		if err != nil {
			return s, err
		}
		str := c.base64EncodeEscape(ciphertext)

		res += str
	}
	return
}

// Decrypt decrypts s3 compatible full filename.
func (c CryptFilename) Decrypt(s string) (res string, err error) {
	parts := strings.Split(s, "/")
	for i, p := range parts {

		// Don't decrypt first folder of path
		if !c.encryptFirst && i == 0 {
			res = p
			continue
		}

		// Path slash
		if i > 0 {
			res += "/"
		}

		// Don't decrypt empy path
		if len(p) == 0 {
			continue
		}

		// Decrypt
		data, err := c.base64DecodeEscape(p)
		if err == nil {
			data, err = crypt.Decrypt(c.hashKey, data)
		}
		if err != nil {
			data = []byte(p)
		}
		data, _ = c.unzip(data)

		res += string(data)
	}
	return
}