package crypt_filename

import (
	"errors"
	"fmt"
	"testing"
)

const key = "very stong key"

func TestCryptFilename(t *testing.T) {

	c := New(key, true)

	fmt.Println()
	fmt.Println("--------------------------------")

	var path, enc, decr string
	var err1, err2 error
	print := func(decrypt_only ...bool) {
		enc = path
		err1 = nil
		if !(len(decrypt_only) > 0 && decrypt_only[0]) {
			enc, err1 = c.Encrypt(path)
		} else {
			fmt.Println("don't encrypt")
		}
		decr, err2 = c.Decrypt(enc)
		fmt.Printf("path: %s\nenc : %s\nerr : %v\ndecr: %s\nerr : %v\n", path, enc, err1, decr, err2)
		if path == decr {
			fmt.Println("done")
		} else {
			fmt.Println("error")
			t.Error(errors.New("input path not equal to decrypted"))
			return
		}
	}

	fmt.Println()
	path = "qqqmmm/path1/path2/file.txt"
	print()

	fmt.Println()
	path = "qqqmmm/long-string-path-long-string-path-long-string-path-long-string-path/long-string-path-long-string-path-long-string-path-long-string-path-file.txt"
	print()

	fmt.Println()
	path = "qqqmmm/"
	print()

	fmt.Println()
	path = "qqqmmm"
	print()

	fmt.Println()
	path = ""
	print()

	fmt.Println()
	path = "qqqmmm/path1/path2/file.txt"
	print(true)

	fmt.Println()
	path = "/qqqmmm/path1/path2/file.txt"
	print()
}
