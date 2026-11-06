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
	//
	// Deprecated: PKCS #1 v1.5 encryption is unsafe. Use EncryptOAEP for new data.
	Encrypt(rawString string, publicKey []byte) (string, error)
	// Decrypt RSA decryption
	//
	// Supported privateKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	//
	// Deprecated: PKCS #1 v1.5 decryption is unsafe. Use DecryptOAEP for OAEP ciphertext.
	Decrypt(encodedString string, privateKey []byte) (string, error)
	// EncryptOAEP encrypts with RSA-OAEP using SHA-256 and the supplied label, returning Base64 ciphertext.
	// It accepts the same public key formats as Encrypt. Plaintext is limited to the key size in bytes minus 66.
	// label may be nil; decryption must use the same label.
	EncryptOAEP(rawString string, publicKey, label []byte) (string, error)
	// DecryptOAEP decrypts Base64 RSA-OAEP ciphertext using SHA-256 and the supplied label.
	// It accepts the same private key formats as Decrypt and does not accept PKCS #1 v1.5 ciphertext.
	// label must match the encryption label; nil and empty labels are equivalent.
	DecryptOAEP(encodedString string, privateKey, label []byte) (string, error)
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
//
// Deprecated: PKCS #1 v1.5 encryption is unsafe. Use EncryptOAEP for new data.
func (rsaImpl) Encrypt(rawString string, publicKey []byte) (string, error) {
	pubObj := parsePublicKey(publicKey)
	if pubObj == nil {
		return "", errors.New("public key is nil")
	}

	encodedByte, err := rsa.EncryptPKCS1v15(rand.Reader, pubObj, []byte(rawString)) //nolint:staticcheck // Preserve PKCS #1 v1.5 ciphertext compatibility for the existing API.
	if err != nil {
		return "", err
	}

	return Base64().Encode(encodedByte), nil
}

// Decrypt RSA decryption
//
// Deprecated: PKCS #1 v1.5 decryption is unsafe. Use DecryptOAEP for OAEP ciphertext.
func (rsaImpl) Decrypt(encodedString string, privateKey []byte) (string, error) {
	encodedByte, err := Base64().Decode(encodedString)
	if err != nil {
		return "", errors.New("invalid base64 encode string")
	}

	priObj := parsePrivateKey(privateKey)
	if priObj == nil {
		return "", errors.New("private key is nil")
	}
	decodedByte, err := rsa.DecryptPKCS1v15(rand.Reader, priObj, encodedByte) //nolint:staticcheck // Existing clients may supply PKCS #1 v1.5 ciphertext; OAEP requires an explicit API migration.
	if err != nil {
		return "", err
	}

	return string(decodedByte), nil
}

// EncryptOAEP encrypts with RSA-OAEP using SHA-256 and the supplied label, returning Base64 ciphertext.
// Plaintext is limited to the key size in bytes minus 66.
// label may be nil; decryption must use the same label.
func (rsaImpl) EncryptOAEP(rawString string, publicKey, label []byte) (string, error) {
	pubObj := parsePublicKey(publicKey)
	if pubObj == nil {
		return "", errors.New("public key is nil")
	}

	encodedByte, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubObj, []byte(rawString), label)
	if err != nil {
		return "", err
	}

	return Base64().Encode(encodedByte), nil
}

// DecryptOAEP decrypts Base64 RSA-OAEP ciphertext using SHA-256 and the supplied label.
// PKCS #1 v1.5 ciphertext is not accepted.
// label must match the encryption label; nil and empty labels are equivalent.
func (rsaImpl) DecryptOAEP(encodedString string, privateKey, label []byte) (string, error) {
	encodedByte, err := Base64().Decode(encodedString)
	if err != nil {
		return "", errors.New("invalid base64 encode string")
	}

	priObj := parsePrivateKey(privateKey)
	if priObj == nil {
		return "", errors.New("private key is nil")
	}

	decodedByte, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priObj, encodedByte, label)
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
	pub, _ := pubInterface.(*rsa.PublicKey)

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
	priObj, _ := pri.(*rsa.PrivateKey)

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
