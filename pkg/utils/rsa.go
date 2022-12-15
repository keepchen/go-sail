package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

//RSAEncrypt rsa加密
func RSAEncrypt(rawString string, publicKey []byte) (string, error) {
	pubObj := parsePublicKey(publicKey)

	encodedByte, err := rsa.EncryptPKCS1v15(rand.Reader, pubObj, []byte(rawString))
	if err != nil {
		return "", err
	}

	return Base64Encode(encodedByte), nil
}

//RSADecrypt rsa解密
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

//RSASign rsa加签
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

//RSAVerifySign rsa验签
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

//解析公钥
func parsePublicKey(key []byte) *rsa.PublicKey {
	block, _ := pem.Decode(key)

	if block == nil {
		return nil
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil
	}
	pub := pubInterface.(*rsa.PublicKey)

	return pub
}

//解析私钥
func parsePrivateKey(key []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil
	}

	pri, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil
	}
	priObj := pri.(*rsa.PrivateKey)

	return priObj
}
