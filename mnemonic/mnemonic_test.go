// Teocryp packge Mnemonic tests
package mnemonic

import (
	"fmt"
	"testing"

	"github.com/teonet-go/teocrypt/config"
)

const appShortName = "teocrypt-test"

func TestNewMnemonic(t *testing.T) {

	// Generate new mnemonic
	mnemonic, err := NewMnemonic()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(mnemonic)
}

func TestNewKeys(t *testing.T) {

	// Generate new mnemonic
	mnemonic, err := NewMnemonic()
	if err != nil {
		t.Error(err)
		return
	}

	// Generate new private and public keys from mnemonic
	privateKey, publicKey, err := GenerateKeys(mnemonic)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("mnemonic: %s\n", mnemonic)
	fmt.Printf("private key: %s\n", privateKey)
	fmt.Printf("public key : %s\n", publicKey)
}

func TestSave(t *testing.T) {

	// Generate new mnemonic
	mnemonic, err := NewMnemonic()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("mnemonic:", mnemonic)

	// Generate private and public keys from mnemonic
	privateKey, _, err := GenerateKeys(string(mnemonic))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("privateKey:", privateKey)

	// Save encrypted mnemonic config
	err = MnemonicConfig{[]byte(mnemonic), []byte(privateKey)}.Save(appShortName)
	if err != nil {
		t.Error(err)
		return
	}

	// Load and print config
	cfg, err := config.Load[MnemonicConfig](appShortName)
	if err != nil {
		t.Error(err)
		return
	}
	data, _ := cfg.Marshal()
	fmt.Println(string(data))
}

func TestLoad(t *testing.T) {
	m := MnemonicConfig{}
	m.Load(appShortName)

	fmt.Println("mnemonic:", string(m.Mnemonic))
	fmt.Println("privateKey:", string(m.PrivateKey))
}
