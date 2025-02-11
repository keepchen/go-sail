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
	// 使用CFB
	//
	// key应该是一个16或24或32位长度的字符
	Encode(rawString, key string) (string, error)
	// Decode aes解密
	//
	// 使用CFB
	//
	// key应该是一个16或24或32位长度的字符
	Decode(encryptedString, key string) (string, error)
}

// Aes 实例化aes工具类
func Aes() IAes {
	return &aesImpl{}
}

var _ IAes = aesImpl{}

// Encode aes加密
//
// 使用CFB
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

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return base64.StdEncoding.EncodeToString(cipherText), nil

}

// Decode aes解密
//
// 使用CFB
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

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
