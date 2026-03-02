package utils

import (
	"encoding/hex"

	"github.com/tjfoc/gmsm/sm4"
)

type sm4impl struct {
}

type ISM4 interface {
	// ECBEncrypt ECB加密
	//
	// hexKey 16进制key 长度32位
	//
	// raw 待加密内容
	ECBEncrypt(hexKey, raw string) (string, error)
	// ECBDecrypt ECB解密
	//
	// hexKey 16进制key 长度32位
	//
	// base64Raw 加密内容 base64格式
	ECBDecrypt(hexKey, base64Raw string) (string, error)
}

var smi ISM4 = &sm4impl{}

// SM4 实例化sm4工具类
func SM4() ISM4 {
	return smi
}

// ECBEncrypt ECB加密
//
// hexKey 16进制key 长度32位
//
// raw 待加密内容
func (sm4impl) ECBEncrypt(hexKey, raw string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}
	out, err := sm4.Sm4Ecb(key, []byte(raw), true)
	if err != nil {
		return "", err
	}

	return Base64().Encode(out), nil
}

// ECBDecrypt ECB解密
//
// hexKey 16进制key 长度32位
//
// base64Raw 加密内容 base64格式
func (sm4impl) ECBDecrypt(hexKey, base64Raw string) (string, error) {
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", err
	}

	raw, err := Base64().Decode(base64Raw)
	if err != nil {
		return "", err
	}

	out, err := sm4.Sm4Ecb(key, raw, false)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
