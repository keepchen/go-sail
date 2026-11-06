package utils

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"
)

type certImpl struct {
}

type ICert interface {
	// ReportValidity reports certificate validity
	//
	// Parameters:
	//
	// - domain: domain name
	//
	// - pemData: certificate chain containing -----BEGIN CERTIFICATE----- xx -----END CERTIFICATE----- in PEM or CRT format
	//
	// Returns:
	//
	// - Whether valid
	//
	// - Certificate validity start time
	//
	// - Certificate validity end time
	//
	// - Error
	ReportValidity(domain string, pemData []byte) (bool, time.Time, time.Time, error)
	// ReportKeyWhetherMatch reports whether the certificate and private key match
	ReportKeyWhetherMatch(certData, keyData []byte) (bool, error)
}

var ci ICert = &certImpl{}

// Cert instantiates certificate utilities
func Cert() ICert {
	return ci
}

// ReportValidity reports certificate validity
//
// Parameters:
//
// - domain: domain name
//
// - pemData: certificate chain containing -----BEGIN CERTIFICATE----- xx -----END CERTIFICATE----- in PEM or CRT format
//
// Returns:
//
// - Whether valid
//
// - Certificate validity start time
//
// - Certificate validity end time
//
// - Error
func (certImpl) ReportValidity(domain string, pemData []byte) (bool, time.Time, time.Time, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return false, time.Now(), time.Now(), fmt.Errorf("the certificate is invalid")
	}
	certObj, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, time.Now(), time.Now(), err
	}
	//Certificate validity ends after it starts
	if certObj.NotAfter.Before(certObj.NotBefore) {
		return false, certObj.NotBefore, certObj.NotAfter, fmt.Errorf("the certificate's validity range is invalid")
	}
	//Certificate validity ends before the current time
	if certObj.NotAfter.Before(time.Now()) {
		return false, certObj.NotBefore, certObj.NotAfter, fmt.Errorf("the certificate has expired")
	}
	//Certificate validity starts after the current time
	if certObj.NotBefore.After(time.Now()) {
		return false, certObj.NotBefore, certObj.NotAfter, fmt.Errorf("the certificate is not yet valid")
	}
	//Check whether the domain matches
	err = certObj.VerifyHostname(domain)

	return err == nil, certObj.NotBefore, certObj.NotAfter, err
}

// ReportKeyWhetherMatch reports whether the certificate and private key match
func (certImpl) ReportKeyWhetherMatch(certData, keyData []byte) (bool, error) {
	certBlock, _ := pem.Decode(certData)
	if certBlock == nil || certBlock.Type != "CERTIFICATE" {
		return false, fmt.Errorf("the certificate is invalid")
	}
	certObj, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return false, err
	}
	certPublicKey, _ := x509.MarshalPKIXPublicKey(certObj.PublicKey)
	keyBlock, _ := pem.Decode(keyData)
	if keyBlock == nil || keyBlock.Type != "PRIVATE KEY" && keyBlock.Type != "RSA PRIVATE KEY" && keyBlock.Type != "EC PRIVATE KEY" {
		return false, fmt.Errorf("the key is invalid")
	}
	var (
		privatePublicKey any
		keyPublicKey     []byte
	)
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err == nil {
		if val, ok := privateKey.(*rsa.PrivateKey); ok {
			privatePublicKey = val.Public()
			keyPublicKey, _ = x509.MarshalPKIXPublicKey(privatePublicKey)
		}
		if val, ok := privateKey.(*ecdsa.PrivateKey); ok {
			privatePublicKey = val.Public()
			keyPublicKey, _ = x509.MarshalPKIXPublicKey(privatePublicKey)
		}
		if val, ok := privateKey.(ed25519.PrivateKey); ok {
			privatePublicKey = val.Public()
			keyPublicKey, _ = x509.MarshalPKIXPublicKey(privatePublicKey)
		}
	} else {
		//Try parsing other private key formats (PKCS1 or EC)
		privateKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
		if err == nil && privateKey != nil {
			privatePublicKey = privateKey.(*rsa.PrivateKey).Public()
			keyPublicKey, _ = x509.MarshalPKIXPublicKey(privatePublicKey)
		}

		if err != nil {
			privateKey, err = x509.ParseECPrivateKey(keyBlock.Bytes)
			if err == nil && privateKey != nil {
				privatePublicKey = privateKey.(ed25519.PrivateKey).Public()
				keyPublicKey, _ = x509.MarshalPKIXPublicKey(privatePublicKey)
			}
		}
	}
	if err != nil {
		return false, fmt.Errorf("the key can not be parse")
	}
	// Compare public keys for equality
	if base64.StdEncoding.EncodeToString(certPublicKey) != base64.StdEncoding.EncodeToString(keyPublicKey) {
		return false, fmt.Errorf("certificate and key not match")
	}

	return true, nil
}
