// Test functions from package arch
package compress

import (
	"bytes"
	"fmt"
	"testing"
)

// TestCompress tests Compress and Decompress function.
func TestCompress(t *testing.T) {
	indata := []byte("Hello, World!, Hello, World!, Hello, World!")

	var in, compressed, decompressed bytes.Buffer

	in.Write(indata)
	err := Compress(&in, &compressed)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	fmt.Println(compressed.Bytes())
	fmt.Println(compressed.String())

	err = Decompress(&compressed, &decompressed)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	fmt.Println(decompressed.String())

	if string(indata) != decompressed.String() {
		t.Errorf("Error: %s", "not equal")
		return
	}
}

// TestCompressData tests CompressData and DecompressData function.
func TestCompressData(t *testing.T) {
	indata := []byte("Hello, World!, Hello, World!, Hello, World!")

	compressed, err := CompressData(indata)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	fmt.Println(compressed)

	decompressed, err := DecompressData(compressed)
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	fmt.Println(string(decompressed))
}
