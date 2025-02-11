package utils

import "encoding/base64"

// Base64Encode base64编码
//
// Deprecated: Base64Encode is deprecated,it will be removed in the future.
//
// Please use Base64().Encode() instead.
func Base64Encode(rawBytes []byte) string {
	return base64.StdEncoding.EncodeToString(rawBytes)
}

// Base64Decode base64解码
//
// Deprecated: Base64Decode is deprecated,it will be removed in the future.
//
// Please use Base64().Decode() instead.
func Base64Decode(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
