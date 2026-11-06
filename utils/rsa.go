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
	// Encrypt RSA encryption
	//
	// Supported publicKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	Encrypt(rawString string, publicKey []byte) (string, error)
	// Decrypt RSA decryption
	//
	// Supported privateKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	Decrypt(encodedString string, privateKey []byte) (string, error)
	// Sign RSA signing
	//
	// Supported privateKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	Sign(rawStringBytes, privateKey []byte) (string, error)
	// VerifySign RSA signature verification
	//
	// Supported publicKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	VerifySign(rawStringBytes, sign, publicKey []byte) (bool, error)
}

var ri IRsa = &rsaImpl{}

// RSA instantiates RSA utilities
func RSA() IRsa {
	return ri
}

// Encrypt RSA encryption
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

// Decrypt RSA decryption
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

// Sign RSA signing
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

// VerifySign RSA signature verification
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

// Parse the public key
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

// Parse the private key
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

// Format public and private keys
//
// Supported formats:
//
// - Standard format
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
// - Nonstandard format: single line without a prefix or suffix
//
// xMIIEvgI...UK2duHjDxiaPKtANi2P4
//
// - Nonstandard format: single line with a prefix or suffix
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
	//Split into lines of 64 characters
	if !strings.Contains(keyStr, "\n") {
		//Remove the prefix first if present
		if strings.HasPrefix(keyStr, prefix) {
			keyStr = strings.Replace(keyStr, prefix, "", 1)
		}
		//Remove the suffix first if present
		if strings.HasSuffix(keyStr, suffix) {
			keyStr = strings.Replace(keyStr, suffix, "", 1)
		}
		keyStr = String().Wordwrap(keyStr, 64, "\n")
	}
	//Check for the prefix and add it if absent
	if !strings.HasPrefix(keyStr, prefix) {
		keyStr = fmt.Sprintf("%s\n%s", prefix, keyStr)
	}
	//Check for the suffix and add it if absent
	if !strings.HasSuffix(keyStr, suffix) {
		keyStr = fmt.Sprintf("%s\n%s", keyStr, suffix)
	}

	return []byte(keyStr)
}
