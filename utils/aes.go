package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type aesImpl struct {
}

// IAes aes接口
type IAes interface {
	// Encode aes加密
	//
	// 使用CTR
	//
	// key应该是一个16或24或32位长度的字符
	Encode(rawString, key string) (string, error)
	// Decode aes解密
	//
	// 使用CTR
	//
	// key应该是一个16或24或32位长度的字符
	Decode(encryptedString, key string) (string, error)
	// GCMEncrypt aes加密
	//
	// 使用GCM
	//
	// key应该是一个16或24或32位长度的字符
	GCMEncrypt(rawString, key string) (string, error)
	// GCMDecrypt aes解密
	//
	// 使用GCM
	//
	// key应该是一个16或24或32位长度的字符
	GCMDecrypt(encryptedString, key string) (string, error)
}

// Aes 实例化aes工具类
func Aes() IAes {
	return ai
}

var ai IAes = &aesImpl{}

// Encode aes加密
//
// 使用CTR
//
// key应该是一个16或24或32位长度的字符
func (aesImpl) Encode(rawString, key string) (string, error) {
	plainText := []byte(rawString)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decode aes解密
//
// 使用CTR
//
// key应该是一个16或24或32位长度的字符
func (aesImpl) Decode(encryptedString, key string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", nil
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("CipherText block size is too short")
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}

// GCMEncrypt aes加密
//
// 使用GCM
//
// key应该是一个16或24或32位长度的字符
func (aesImpl) GCMEncrypt(rawString, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
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
	cipherText := gcm.Seal(nonce, nonce, []byte(rawString), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// GCMDecrypt aes解密
//
// 使用GCM
//
// key应该是一个16或24或32位长度的字符
func (aesImpl) GCMDecrypt(encryptedString, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext is too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
