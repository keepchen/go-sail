package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

type gzipImpl struct {
}

type IGzip interface {
	// Compress gzip compression
	//
	// Uses the gzip.DefaultCompression level
	Compress(content []byte) []byte
	// Decompress gzip decompression
	Decompress(content []byte) []byte
}

var gi IGzip = &gzipImpl{}

// Gzip instantiates gzip utilities
func Gzip() IGzip {
	return gi
}

// Compress gzip compression
//
// Uses the gzip.DefaultCompression level
func (gzipImpl) Compress(content []byte) []byte {
	var b bytes.Buffer
	w, e := gzip.NewWriterLevel(&b, gzip.DefaultCompression)
	if e != nil {
		return content
	}

	_, _ = w.Write(content)
	_ = w.Close()

	return b.Bytes()
}

// Decompress gzip decompression
func (gzipImpl) Decompress(content []byte) []byte {
	r, e := gzip.NewReader(bytes.NewReader(content))
	if e != nil {
		return []byte(``)
	}
	res, _ := io.ReadAll(r)
	_ = r.Close()

	return res
}
