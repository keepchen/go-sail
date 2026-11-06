package sail

import "github.com/keepchen/go-sail/v3/utils"

type passwordImpl struct{}

// IPassword password utility interface
type IPassword interface {
	// AesGCMEncrypt AES encryption
	//
	// Uses GCM
	//
	// key should be a string of length 16, 24, or 32
	AesGCMEncrypt(plaintext, key string) (string, error)
	// AesGCMDecrypt AES decryption
	//
	// Uses GCM
	//
	// key should be a string of length 16, 24, or 32
	AesGCMDecrypt(ciphertext, key string) (string, error)
	// SM4GCMEncrypt GCM encryption
	//
	// hexKey hexadecimal key, 32 characters long
	//
	// raw content to encrypt
	SM4GCMEncrypt(plaintext, key string) (string, error)
	// SM4GCMDecrypt GCM decryption
	//
	// hexKey hexadecimal key, 32 characters long
	//
	// base64Raw encrypted content in Base64 format
	SM4GCMDecrypt(ciphertext, key string) (string, error)
	// RSAEncrypt RSA encryption
	//
	// Supported publicKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	//
	// Deprecated: PKCS #1 v1.5 encryption is unsafe. Use RSAEncryptOAEP for new data.
	RSAEncrypt(plaintext string, publicKey []byte) (string, error)
	// RSADecrypt RSA decryption
	//
	// Supported privateKey formats:
	//
	// - Standard format
	//
	// - Single-line string without the prefix and suffix
	//
	// Deprecated: PKCS #1 v1.5 decryption is unsafe. Use RSADecryptOAEP for OAEP ciphertext.
	RSADecrypt(ciphertext string, privateKey []byte) (string, error)
	// RSAEncryptOAEP encrypts with RSA-OAEP using SHA-256 and an empty label, returning Base64 ciphertext.
	// It accepts the same public key formats as RSAEncrypt. Plaintext is limited to the key size in bytes minus 66.
	RSAEncryptOAEP(plaintext string, publicKey []byte) (string, error)
	// RSADecryptOAEP decrypts Base64 RSA-OAEP ciphertext using SHA-256 and an empty label.
	// It accepts the same private key formats as RSADecrypt and does not accept PKCS #1 v1.5 ciphertext.
	RSADecryptOAEP(ciphertext string, privateKey []byte) (string, error)
	// BcryptHashMake generates a hash string
	//
	// plaintext - plaintext string
	//
	// costs - cost factor, range ∈ [4,31], default: 10
	BcryptHashMake(plaintext string, costs ...int) (string, error)
	// BcryptHashCheck verifies a hash string
	//
	// plaintext - plaintext string
	//
	// hashed - hashed string
	BcryptHashCheck(plaintext string, hash string) (bool, error)
	// Argon2HashMake generates a hash string
	//
	// plaintext - plaintext string
	//
	// argon2Config - configuration parameters
	Argon2HashMake(plaintext string, argon2Config ...utils.Argon2Config) (string, error)
	// Argon2HashCheck verifies a hash string
	//
	// plaintext - plaintext string
	//
	// hashed - hashed string
	//
	// argon2Config - configuration parameters
	Argon2HashCheck(plaintext string, hash string, argon2Config ...utils.Argon2Config) (bool, error)
}

var _ IPassword = &passwordImpl{}

// Password instantiates password utilities
func Password() IPassword {
	return &passwordImpl{}
}

func (*passwordImpl) AesGCMEncrypt(plaintext, key string) (string, error) {
	return utils.Aes().GCMEncrypt(plaintext, key)
}

func (*passwordImpl) AesGCMDecrypt(ciphertext, key string) (string, error) {
	return utils.Aes().GCMDecrypt(ciphertext, key)
}

func (*passwordImpl) SM4GCMEncrypt(plaintext, key string) (string, error) {
	return utils.SM4().GCMEncrypt(key, plaintext)
}

func (*passwordImpl) SM4GCMDecrypt(ciphertext, key string) (string, error) {
	return utils.SM4().GCMDecrypt(key, ciphertext)
}

// RSAEncrypt encrypts with RSA PKCS #1 v1.5.
//
// Deprecated: PKCS #1 v1.5 encryption is unsafe. Use RSAEncryptOAEP for new data.
func (*passwordImpl) RSAEncrypt(plaintext string, publicKey []byte) (string, error) {
	return utils.RSA().Encrypt(plaintext, publicKey)
}

// RSADecrypt decrypts RSA PKCS #1 v1.5 ciphertext.
//
// Deprecated: PKCS #1 v1.5 decryption is unsafe. Use RSADecryptOAEP for OAEP ciphertext.
func (*passwordImpl) RSADecrypt(ciphertext string, privateKey []byte) (string, error) {
	return utils.RSA().Decrypt(ciphertext, privateKey)
}

func (*passwordImpl) RSAEncryptOAEP(plaintext string, publicKey []byte) (string, error) {
	return utils.RSA().EncryptOAEP(plaintext, publicKey, nil)
}

func (*passwordImpl) RSADecryptOAEP(ciphertext string, privateKey []byte) (string, error) {
	return utils.RSA().DecryptOAEP(ciphertext, privateKey, nil)
}

func (*passwordImpl) BcryptHashMake(plaintext string, costs ...int) (string, error) {
	return utils.Bcrypt().HashMake(plaintext, costs...)
}

func (*passwordImpl) BcryptHashCheck(plaintext string, hash string) (bool, error) {
	return utils.Bcrypt().HashCheck(plaintext, hash)
}

func (*passwordImpl) Argon2HashMake(plaintext string, argon2Config ...utils.Argon2Config) (string, error) {
	return utils.Argon2(argon2Config...).HashMake(plaintext)
}

func (*passwordImpl) Argon2HashCheck(plaintext string, hash string, argon2Config ...utils.Argon2Config) (bool, error) {
	return utils.Argon2(argon2Config...).HashCheck(plaintext, hash)
}
