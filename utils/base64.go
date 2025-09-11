package utils

import "encoding/base64"

type base64Impl struct {
}

// IBase64 base64接口
type IBase64 interface {
	// Encode base64编码
	Encode(rawBytes []byte) string
	// Decode base64解码
	Decode(encodedString string) ([]byte, error)
}

// Base64 实例化base64工具类
func Base64() IBase64 {
	return bi
}

var bi IBase64 = &base64Impl{}

// Encode base64编码
func (base64Impl) Encode(rawBytes []byte) string {
	return base64.StdEncoding.EncodeToString(rawBytes)
}

// Decode base64解码
func (base64Impl) Decode(encodedString string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encodedString)
}
