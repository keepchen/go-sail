package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/keepchen/go-sail/v3/constants"
)

type rsaImpl struct {
}

type IRsa interface {
	// Encrypt rsa加密
	//
	// publicKey格式支持：
	//
	// - 标准格式
	//
	// - 不带前后缀的一行字符串
	Encrypt(rawString string, publicKey []byte) (string, error)
	// Decrypt rsa解密
	//
	// privateKey格式支持：
	//
	// - 标准格式
	//
	// - 不带前后缀的一行字符串
	Decrypt(encodedString string, privateKey []byte) (string, error)
	// Sign rsa加签
	//
	// privateKey格式支持：
	//
	// - 标准格式
	//
	// - 不带前后缀的一行字符串
	Sign(rawStringBytes, privateKey []byte) (string, error)
	// VerifySign rsa验签
	//
	// publicKey格式支持：
	//
	// - 标准格式
	//
	// - 不带前后缀的一行字符串
	VerifySign(rawStringBytes, sign, publicKey []byte) (bool, error)
}

// RSA 实例化rsa工具类
func RSA() IRsa {
	return &rsaImpl{}
}

// Encrypt rsa加密
func (rsaImpl) Encrypt(rawString string, publicKey []byte) (string, error) {
	pubObj := parsePublicKey(publicKey)
	if pubObj == nil {
		return "", errors.New("public key is nil")
	}

	encodedByte, err := rsa.EncryptPKCS1v15(rand.Reader, pubObj, []byte(rawString))
	if err != nil {
		return "", err
	}

	return Base64().Encode(encodedByte), nil
}

// Decrypt rsa解密
func (rsaImpl) Decrypt(encodedString string, privateKey []byte) (string, error) {
	encodedByte, err := Base64().Decode(encodedString)
	if err != nil {
		return "", errors.New("invalid base64 encode string")
	}

	priObj := parsePrivateKey(privateKey)
	if priObj == nil {
		return "", errors.New("private key is nil")
	}
	decodedByte, err := rsa.DecryptPKCS1v15(rand.Reader, priObj, encodedByte)
	if err != nil {
		return "", err
	}

	return string(decodedByte), nil
}

// Sign rsa加签
func (rsaImpl) Sign(rawStringBytes, privateKey []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(rawStringBytes)
	if err != nil {
		return "", err
	}
	d := h.Sum(nil)

	priObj := parsePrivateKey(privateKey)
	if priObj == nil {
		return "", errors.New("private key is nil")
	}

	sign, err := rsa.SignPKCS1v15(rand.Reader, priObj, crypto.SHA256, d)

	return Base64().Encode(sign), err
}

// VerifySign rsa验签
func (rsaImpl) VerifySign(rawStringBytes, sign, publicKey []byte) (bool, error) {
	h := sha256.New()
	_, err := h.Write(rawStringBytes)
	if err != nil {
		return false, err
	}
	d := h.Sum(nil)

	pubObj := parsePublicKey(publicKey)
	if pubObj == nil {
		return false, errors.New("public key is nil")
	}

	decodedSign, err := Base64().Decode(string(sign))
	if err != nil {
		return false, err
	}

	err = rsa.VerifyPKCS1v15(pubObj, crypto.SHA256, d, decodedSign)

	return err == nil, err
}

// 解析公钥
func parsePublicKey(key []byte) *rsa.PublicKey {
	key = formatKey(key, "public")
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

// 解析私钥
func parsePrivateKey(key []byte) *rsa.PrivateKey {
	key = formatKey(key, "private")
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

// 格式化公私钥
//
// 支持格式：
//
// - 标准格式
//
// -----BEGIN PRIVATE KEY-----
//
// MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDUvUDx+LPQ0S+L
//
// ...
//
// xUK2duHjDxiaPKtANi2P4
//
// -----END PRIVATE KEY-----
//
// - 非标准格式：一行不带前缀或不带后缀
//
// xMIIEvgI...UK2duHjDxiaPKtANi2P4
//
// - 非标准格式：一行带前缀或后缀
//
// -----BEGIN PRIVATE KEY-----MIIEvgI...UK2duHjDxiaPKtANi2P4-----END PRIVATE KEY-----
func formatKey(key []byte, typ string) []byte {
	var (
		prefix string
		suffix string
		keyStr = string(key)
	)
	if typ == "public" {
		prefix = constants.PublicKeyBeginStr
		suffix = constants.PublicKeyEndStr
	} else {
		prefix = constants.PrivateKeyBeginStr
		suffix = constants.PrivateKeyEndStr
	}
	//按每行64个字符拆分
	if !strings.Contains(keyStr, "\n") {
		//存在前缀，先去掉
		if strings.HasPrefix(keyStr, prefix) {
			keyStr = strings.Replace(keyStr, prefix, "", 1)
		}
		//存在后缀，先去掉
		if strings.HasSuffix(keyStr, suffix) {
			keyStr = strings.Replace(keyStr, suffix, "", 1)
		}
		keyStr = String().Wordwrap(keyStr, 64, "\n")
	}
	//检测前缀，没有则添加上
	if !strings.HasPrefix(keyStr, prefix) {
		keyStr = fmt.Sprintf("%s\n%s", prefix, keyStr)
	}
	//检测后缀，没有则添加上
	if !strings.HasSuffix(keyStr, suffix) {
		keyStr = fmt.Sprintf("%s\n%s", keyStr, suffix)
	}

	return []byte(keyStr)
}
