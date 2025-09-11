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
	// ReportValidity 报告证书有效性
	//
	// 参数：
	//
	// - domain: 域名
	//
	// - pemData: 包含 -----BEGIN CERTIFICATE----- xx -----END CERTIFICATE----- 的证书链，格式为pem或crt
	//
	// 返回值：
	//
	// - 是否有效
	//
	// - 证书有效期开始时间
	//
	// - 证书有效期结束时间
	//
	// - 错误
	ReportValidity(domain string, pemData []byte) (bool, time.Time, time.Time, error)
	// ReportKeyWhetherMatch 报告证书和私钥是否匹配
	ReportKeyWhetherMatch(certData, keyData []byte) (bool, error)
}

var ci ICert = &certImpl{}

// Cert 实例化证书工具类
func Cert() ICert {
	return ci
}

// ReportValidity 报告证书有效性
//
// 参数：
//
// - domain: 域名
//
// - pemData: 包含 -----BEGIN CERTIFICATE----- xx -----END CERTIFICATE----- 的证书链，格式为pem或crt
//
// 返回值：
//
// - 是否有效
//
// - 证书有效期开始时间
//
// - 证书有效期结束时间
//
// - 错误
func (certImpl) ReportValidity(domain string, pemData []byte) (bool, time.Time, time.Time, error) {
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "CERTIFICATE" {
		return false, time.Now(), time.Now(), fmt.Errorf("the certificate is invalid")
	}
	certObj, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, time.Now(), time.Now(), err
	}
	//证书有效期结束时间晚于开始时间
	if certObj.NotAfter.Before(certObj.NotBefore) {
		return false, certObj.NotBefore, certObj.NotAfter, fmt.Errorf("the certificate's validity range is invalid")
	}
	//证书有效期结束时间早于当前时间
	if certObj.NotAfter.Before(time.Now()) {
		return false, certObj.NotBefore, certObj.NotAfter, fmt.Errorf("the certificate has expired")
	}
	//证书有效期开始时间晚于当前时间
	if certObj.NotBefore.After(time.Now()) {
		return false, certObj.NotBefore, certObj.NotAfter, fmt.Errorf("the certificate is not yet valid")
	}
	//验证域名是否匹配
	err = certObj.VerifyHostname(domain)

	return err == nil, certObj.NotBefore, certObj.NotAfter, err
}

// ReportKeyWhetherMatch 报告证书和私钥是否匹配
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
		privatePublicKey interface{}
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
		//尝试解析其他格式的私钥（PKCS1 或 EC）
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
	// 比较公钥是否一致
	if base64.StdEncoding.EncodeToString(certPublicKey) != base64.StdEncoding.EncodeToString(keyPublicKey) {
		return false, fmt.Errorf("certificate and key not match")
	}

	return true, nil
}
