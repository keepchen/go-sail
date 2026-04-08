package utils

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"

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
	// GCMEncrypt GCM加密
	//
	// hexKey 16进制key 长度16位
	//
	// raw 待加密内容
	GCMEncrypt(hexKey, raw string) (string, error)
	// GCMDecrypt GCM解密
	//
	// hexKey 16进制key 长度16位
	//
	// base64Raw 加密内容 base64格式
	GCMDecrypt(hexKey, base64Raw string) (string, error)
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

// GCMEncrypt GCM加密
//
// hexKey 16进制key 长度32位
//
// raw 待加密内容
func (sm4impl) GCMEncrypt(hexKey, raw string) (string, error) {
	block, err := sm4.NewCipher([]byte(hexKey))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(raw), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// GCMDecrypt GCM解密
//
// hexKey 16进制key 长度32位
//
// base64Raw 加密内容 base64格式
func (sm4impl) GCMDecrypt(hexKey, base64Raw string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(base64Raw)
	if err != nil {
		return "", err
	}
	block, err := sm4.NewCipher([]byte(hexKey))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
