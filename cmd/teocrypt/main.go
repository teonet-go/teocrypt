// Copyright 2021 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The Teocrypt application is used to encrypt and decrypt text using a key or
// password on the command line.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"

	"github.com/teonet-go/teocrypt/crypt"
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
		key = crypt.HashKey(passwd)
	} else if key, err = crypt.GenerateKey(); err != nil {
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
		ciphertext, err := crypt.Encrypt(key, data)
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
	plaintext, err := crypt.Decrypt(key, data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("plain text:\n%s\n", plaintext)
}
