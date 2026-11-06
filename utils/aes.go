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

// IAes AES interface
type IAes interface {
	// Encode AES encryption
	//
	// Uses CTR
	//
	// key should be a string of length 16, 24, or 32
	Encode(rawString, key string) (string, error)
	// Decode AES decryption
	//
	// Uses CTR
	//
	// key should be a string of length 16, 24, or 32
	Decode(encryptedString, key string) (string, error)
	// GCMEncrypt AES encryption
	//
	// Uses GCM
	//
	// key should be a string of length 16, 24, or 32
	GCMEncrypt(rawString, key string) (string, error)
	// GCMDecrypt AES decryption
	//
	// Uses GCM
	//
	// key should be a string of length 16, 24, or 32
	GCMDecrypt(encryptedString, key string) (string, error)
}

// Aes instantiates AES utilities
func Aes() IAes {
	return ai
}

var ai IAes = &aesImpl{}

// Encode AES encryption
//
// Uses CTR
//
// key should be a string of length 16, 24, or 32
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

// Decode AES decryption
//
// Uses CTR
//
// key should be a string of length 16, 24, or 32
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

// GCMEncrypt AES encryption
//
// Uses GCM
//
// key should be a string of length 16, 24, or 32
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

// GCMDecrypt AES decryption
//
// Uses GCM
//
// key should be a string of length 16, 24, or 32
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
