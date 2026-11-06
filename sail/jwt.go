package sail

import (
	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/keepchen/go-sail/v3/utils"
)

// IJwt JWT interface
type IJwt interface {
	// MakeToken generates a JWT token
	//
	// uid user ID
	//
	// exp token expiration, Unix timestamp in seconds
	//
	// otherFields additional extensible fields
	MakeToken(uid string, exp int64, otherFields ...map[string]any) (string, error)
	// ValidToken validates a JWT token
	//
	// # ⚠️Note⚠️
	//
	// This method inherits jwt.Validator's Validate method, so it only validates fields that are present and correctly formatted;
	// callers must implement their own validation logic for any additional rules.
	ValidToken(token string) (bool, jwt.MapClaims, error)
	//Encrypt encrypts a string
	//
	// This method pairs with Decrypt for encryption and decryption
	//
	// # Note
	//
	// This method calls utils.RSA and encrypts with PublicKey from jwt.Conf
	//
	// Deprecated: PKCS #1 v1.5 encryption is unsafe. Use EncryptOAEP for new data.
	Encrypt(plaintext string) (string, error)
	//Decrypt decrypts a string
	//
	// This method pairs with Encrypt for encryption and decryption
	//
	// # Note
	//
	// This method calls utils.RSA and decrypts with PrivateKey from jwt.Conf
	//
	// Deprecated: PKCS #1 v1.5 decryption is unsafe. Use DecryptOAEP for OAEP ciphertext.
	Decrypt(encryptedStr string) (string, error)
	// EncryptOAEP encrypts with PublicKey from jwt.Conf using RSA-OAEP, SHA-256, and an empty label.
	// It returns Base64 ciphertext. Plaintext is limited to the key size in bytes minus 66.
	EncryptOAEP(plaintext string) (string, error)
	// DecryptOAEP decrypts Base64 RSA-OAEP ciphertext with PrivateKey from jwt.Conf using SHA-256 and an empty label.
	// PKCS #1 v1.5 ciphertext is not accepted.
	DecryptOAEP(encryptedStr string) (string, error)
}

type jwtImpl struct{}

var _ IJwt = &jwtImpl{}

// JWT JWT instance
func JWT() IJwt {
	return &jwtImpl{}
}

// MakeToken generates a JWT token
//
// uid user ID
//
// exp token expiration, Unix timestamp in seconds
//
// otherFields additional extensible fields
func (jwtImpl) MakeToken(uid string, exp int64, otherFields ...map[string]any) (string, error) {
	conf := config.Get()
	baseMap := map[string]any{
		"uid": uid,
		"iss": conf.JwtConf.TokenIssuer,
		"exp": exp,
	}
	//Values may be overwritten
	for _, otherField := range otherFields {
		for k, v := range otherField {
			baseMap[k] = v
		}
	}
	mp := jwt.MergeStandardClaims(baseMap)

	return jwt.SignWithMap(mp, *conf.JwtConf)
}

// ValidToken validates a JWT token
//
// # ⚠️Note⚠️
//
// This method inherits jwt.Validator's Validate method, so it only validates fields that are present and correctly formatted;
// callers must implement their own validation logic for any additional rules.
func (jwtImpl) ValidToken(token string) (bool, jwt.MapClaims, error) {
	conf := config.Get()
	mp, err := jwt.VerifyFromMap(token, *conf.JwtConf)

	return err == nil, mp, err
}

// Encrypt encrypts a string
//
// # This method pairs with Decrypt for encryption and decryption
//
// # Note
//
// # This method calls utils.RSA and encrypts with PublicKey from jwt.Conf
//
// Deprecated: PKCS #1 v1.5 encryption is unsafe. Use EncryptOAEP for new data.
func (jwtImpl) Encrypt(plaintext string) (string, error) {
	conf := config.Get()
	return utils.RSA().Encrypt(plaintext, []byte(conf.JwtConf.PublicKey))
}

// Decrypt decrypts a string
//
// # This method pairs with Encrypt for encryption and decryption
//
// # Note
//
// # This method calls utils.RSA and decrypts with PrivateKey from jwt.Conf
//
// Deprecated: PKCS #1 v1.5 decryption is unsafe. Use DecryptOAEP for OAEP ciphertext.
func (jwtImpl) Decrypt(encryptedStr string) (string, error) {
	conf := config.Get()
	return utils.RSA().Decrypt(encryptedStr, []byte(conf.JwtConf.PrivateKey))
}

// EncryptOAEP encrypts with PublicKey from jwt.Conf using RSA-OAEP, SHA-256, and an empty label.
func (jwtImpl) EncryptOAEP(plaintext string) (string, error) {
	conf := config.Get()
	return utils.RSA().EncryptOAEP(plaintext, []byte(conf.JwtConf.PublicKey), nil)
}

// DecryptOAEP decrypts Base64 RSA-OAEP ciphertext with PrivateKey from jwt.Conf using SHA-256 and an empty label.
func (jwtImpl) DecryptOAEP(encryptedStr string) (string, error) {
	conf := config.Get()
	return utils.RSA().DecryptOAEP(encryptedStr, []byte(conf.JwtConf.PrivateKey), nil)
}
