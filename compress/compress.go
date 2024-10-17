// Copyright 2024 Kirill Scherba <kirill@scherba.ru>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Archive package contains functions to archive and unarchive strings, data
// and files.
package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

// Compress writes gzipped data from a reader to a writer.
//
// It uses the default compression level.
func Compress(r io.Reader, w io.Writer) (err error) {
	// Create a gzip writer
	gw, err := gzip.NewWriterLevel(w, gzip.DefaultCompression)
	if err != nil {
		return
	}
	defer gw.Close()

	// Copy the data from the reader to the writer
	_, err = io.Copy(gw, r)
	return
}

// Decompress writes gunzipped data from a reader to a writer.
//
// It reads gziped data from the reader and writes it to the writer.
func Decompress(r io.Reader, w io.Writer) (err error) {
	// Read gziped data from the reader
	gr, err := gzip.NewReader(r)
	if err != nil {
		return
	}
	defer gr.Close()

	// Copy the data from the gzip reader to the writer
	_, err = io.Copy(w, gr)
	return
}

// CompressData compresses data.
func CompressData(data []byte) (compressed []byte, err error) {

	var r, w bytes.Buffer

	r.Write(data)
	err = Compress(&r, &w)
	if err != nil {
		return
	}

	compressed = w.Bytes()
	return
}

// DecompressData decompresses data.
func DecompressData(compressed []byte) (decompressed []byte, err error) {

	var r, w bytes.Buffer

	r.Write(compressed)
	err = Decompress(&r, &w)
	if err != nil {
		return
	}

	decompressed = w.Bytes()
	return
}
