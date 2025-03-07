package utils

import (
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm4"
)

// SM4ECBEncrypt ECB加密
//
// Deprecated: SM4ECBEncrypt is deprecated,it will be removed in the future.
//
// Please use SM4().ECBEncrypt() instead.
//
// hexKey 16进制key 长度32位
//
// raw 待加密内容
func SM4ECBEncrypt(hexKey, raw string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}
	out, err := sm4.Sm4Ecb(key, []byte(raw), true)
	if err != nil {
		return "", err
	}

	return Base64Encode(out), nil
}

// SM4ECBDecrypt ECB解密
//
// Deprecated: SM4ECBDecrypt is deprecated,it will be removed in the future.
//
// Please use SM4().ECBDecrypt() instead.
//
// hexKey 16进制key 长度32位
//
// base64Raw 加密内容 base64格式
func SM4ECBDecrypt(hexKey, base64Raw string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}

	raw, err := Base64Decode(base64Raw)
	if err != nil {
		return "", err
	}

	out, err := sm4.Sm4Ecb(key, raw, false)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
