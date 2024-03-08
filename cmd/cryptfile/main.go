// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The Cryptfile application is used to encrypt and decrypt file using a key or
// password on the command line.
//
//	Usage:
//	# Encrypt file:
//	  go run ./cmd/cryptfile/ -p 123456 -i go.mod -o go.mod.crypt
//	  cat go.mod.crypt | base64
//	# Decrypt file:
//	  go run ./cmd/cryptfile/ -p 123456 -d -i go.mod.crypt -o go.mod.decrypt
//	  cat go.mod.decrypt
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/teonet-go/teocrypt/crypt"
)

const (
	appShort   = "cryptfile"
	appName    = "Teonet encrypt/decrypt file application"
	appLong    = ""
	appVersion = "0.0.1"
)

func main() {
	// Application logo
	fmt.Println(appName + " ver " + appVersion)

	// Parse application command line parameters
	var save, decrypt bool
	var inFile, outFile, passwd string
	flag.StringVar(&inFile, "i", "", "input file to encrypt/decrypt")
	flag.StringVar(&outFile, "o", "", "output file to encrypt/decrypt")
	flag.StringVar(&passwd, "p", "", "password used to encrypt/decrypt")
	flag.BoolVar(&decrypt, "d", decrypt, "decrypt file specified in -i flag")
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

	// Open input and output files
	inputFile, _ := os.Open(inFile)
	defer inputFile.Close()
	outputFile, _ := os.Create(outFile)
	defer outputFile.Close()

	if !decrypt {

		// Create Encrypt writer using outputFile writer and key
		var writer io.Writer
		writer, err = crypt.EncryptWriter(outputFile, key)
		if err != nil {
			fmt.Printf("can't create encrypt writer, error: %s\n", err)
			return
		}

		// Copy data from input to output using the stream cipher writer
		_, err = io.Copy(writer, inputFile)

	} else {

		// Create Dencrypt reader using inputFile reader and key
		var reader io.Reader
		reader, err = crypt.DecryptReader(inputFile, key)
		if err != nil {
			fmt.Printf("can't create encrypt writer, error: %s\n", err)
			return
		}

		// Copy data from input to output using the stream cipher reader
		_, err = io.Copy(outputFile, reader)
	}

	if err != nil {
		fmt.Printf("can't execute command, error: %s\n", err)
		return
	}

	fmt.Println("done")
}
