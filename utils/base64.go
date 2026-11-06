package utils

import "encoding/base64"

type base64Impl struct {
}

// IBase64 Base64 interface
type IBase64 interface {
	// Encode Base64 encoding
	Encode(rawBytes []byte) string
	// Decode Base64 decoding
	Decode(encodedString string) ([]byte, error)
}

// Base64 instantiates Base64 utilities
func Base64() IBase64 {
	return bi
}

var bi IBase64 = &base64Impl{}

// Encode Base64 encoding
func (base64Impl) Encode(rawBytes []byte) string {
	return base64.StdEncoding.EncodeToString(rawBytes)
}

// Decode Base64 decoding
func (base64Impl) Decode(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
