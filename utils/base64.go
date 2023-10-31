package utils

import "encoding/base64"

//Base64Encode base64编码
func Base64Encode(rawBytes []byte) string {
	return base64.StdEncoding.EncodeToString(rawBytes)
}

//Base64Decode base64解码
func Base64Decode(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
