// Copyright 2023 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Teocrypt package contains Mnemonic encription functions.
package mnemonic

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/teonet-go/teocrypt/config"
	"github.com/teonet-go/teocrypt/crypt"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// MnemonicConfig is data used in config file
type MnemonicConfig struct {
	Mnemonic   []byte `json:"mnemonic"`
	PrivateKey []byte `json:"private_key"`
}

// NewMnemonic generates a mnemonic string.
func NewMnemonic() (mnemonic string, err error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return
	}
	mnemonic, err = bip39.NewMnemonic(entropy)
	return
}

// GenerateKeys generates master private and public keys from mnemonic.
func GenerateKeys(mnemonic string) (privateKey string, publicKey string, err error) {

	// Generate a Bip32 HD wallet from the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, "Secret Passphrase")

	// Create master private key from seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return
	}

	// Return keys
	privateKey = masterKey.B58Serialize()
	publicKey = masterKey.PublicKey().B58Serialize()
	return
}

// generateChildPrivateKey generates a child private key from a master private key by
// index.
func generateChildPrivateKey(masterKey *bip32.Key, childIdx uint32) (key *bip32.Key, err error) {
	key, err = masterKey.NewChildKey(childIdx)
	if err != nil {
		return
	}
	return
}

// GenerateChildKey generates a child public key from a master private key by
// index.
func GenerateChildKey(masterKeyB58 string, childIdx uint32) (key string, err error) {

	// Deserialize master private key
	masterKey, err := bip32.B58Deserialize(masterKeyB58)
	if err != nil {
		return
	}

	// Generate child private key from master private key
	childPrivateKey, err := generateChildPrivateKey(masterKey, childIdx)
	if err != nil {
		return
	}

	// Get child public key from child private key and serialize it
	key = childPrivateKey.PublicKey().B58Serialize()

	return
}

// Save saves encrypted by "machineid + password" mnemonic config to this
// machine on os.UserConfig/teonet_config_dir/appShortName folder.
func (m MnemonicConfig) Save(appShortName, configName string, passwd ...string) (err error) {

	// Get machine id key to encrypt mnemonic config
	key, err := getKey(passwd...)
	if err != nil {
		return
	}

	// Encrypt mnemonic
	m.Mnemonic, err = crypt.Encrypt(key, m.Mnemonic)
	if err != nil {
		return err
	}

	// Encrypt private key
	m.PrivateKey, err = crypt.Encrypt(key, m.PrivateKey)
	if err != nil {
		return err
	}

	// Save config
	cfg, err := config.New[MnemonicConfig](appShortName, configName, &m)
	if err != nil {
		return
	}
	cfg.Save()

	return
}

// Load loads from config file and decrypt saved mnemonic config.
func (m *MnemonicConfig) Load(appShortName, configName string, passwd ...string) (err error) {

	// Load config
	_, err = config.Load[MnemonicConfig](appShortName, configName, m)
	if err != nil {
		return
	}

	// Get key to decrypt mnemonic
	key, err := getKey()
	if err != nil {
		return
	}

	// Decrypt mnemonic
	m.Mnemonic, err = crypt.Decrypt(key, m.Mnemonic)
	if err != nil {
		return
	}

	// Decrypt private key
	m.PrivateKey, err = crypt.Decrypt(key, m.PrivateKey)
	if err != nil {
		return
	}

	return
}

// IsMnemonicValid attempts to verify that the provided mnemonic is valid.
// Validity is determined by both the number of words being appropriate,
// and that all the words in the mnemonic are present in the word list.
func IsMnemonicValid(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// getKey generates and returns key created from "machineid + password"
func getKey(passwd ...string) (key []byte, err error) {

	// Get unical machine id
	id, err := machineid.ID()
	if err != nil {
		return
	}

	// Add password to the machine id key
	if len(passwd) > 0 {
		id += passwd[0]
	}

	// Generate hash
	key = crypt.HashKey(id)

	return
}
