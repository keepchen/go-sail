package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

// RSAEncrypt rsa加密
//
// Deprecated: RSAEncrypt is deprecated,it will be removed in the future.
//
// Please use RSA().Encrypt() instead.
func RSAEncrypt(rawString string, publicKey []byte) (string, error) {
	pubObj := parsePublicKey(publicKey)

	encodedByte, err := rsa.EncryptPKCS1v15(rand.Reader, pubObj, []byte(rawString))
	if err != nil {
		return "", err
	}

	return Base64Encode(encodedByte), nil
}

// RSADecrypt rsa解密
//
// Deprecated: RSADecrypt is deprecated,it will be removed in the future.
//
// Please use RSA().Decrypt() instead.
func RSADecrypt(encodedString string, privateKey []byte) (string, error) {
	encodedByte, err := Base64Decode(encodedString)
	if err != nil {
		return "", errors.New("invalid base64 encode string")
	}

	priObj := parsePrivateKey(privateKey)
	decodedByte, err := rsa.DecryptPKCS1v15(rand.Reader, priObj, encodedByte)
	if err != nil {
		return "", err
	}

	return string(decodedByte), nil
}

// RSASign rsa加签
//
// Deprecated: RSASign is deprecated,it will be removed in the future.
//
// Please use RSA().Sign() instead.
func RSASign(rawStringBytes, privateKey []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(rawStringBytes)
	if err != nil {
		return "", err
	}
	d := h.Sum(nil)

	priObj := parsePrivateKey(privateKey)

	sign, err := rsa.SignPKCS1v15(rand.Reader, priObj, crypto.SHA256, d)

	return Base64Encode(sign), err
}

// RSAVerifySign rsa验签
//
// Deprecated: RSAVerifySign is deprecated,it will be removed in the future.
//
// Please use RSA().VerifySign() instead.
func RSAVerifySign(rawStringBytes, sign, publicKey []byte) (bool, error) {
	h := sha256.New()
	_, err := h.Write(rawStringBytes)
	if err != nil {
		return false, err
	}
	d := h.Sum(nil)

	pubObj := parsePublicKey(publicKey)

	decodedSign, err := Base64Decode(string(sign))
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(pubObj, crypto.SHA256, d, decodedSign)

	return err == nil, err
}
