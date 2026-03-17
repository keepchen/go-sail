package sail

import "github.com/keepchen/go-sail/v3/utils"

type passwordImpl struct{}

// IPassword 密码工具类接口
type IPassword interface {
	AesGCMEncrypt(plaintext, key string) (string, error)
	AesGCMDecrypt(ciphertext, key string) (string, error)
	SM4GCMEncrypt(plaintext, key string) (string, error)
	SM4GCMDecrypt(ciphertext, key string) (string, error)
	RSAEncrypt(plaintext string, publicKey []byte) (string, error)
	RSADecrypt(ciphertext string, privateKey []byte) (string, error)
	BcryptHashMake(plaintext string, costs ...int) (string, error)
	BcryptHashCheck(plaintext string, hash string) (bool, error)
	Argon2HashMake(plaintext string, argon2Config ...utils.Argon2Config) (string, error)
	Argon2HashCheck(plaintext string, hash string, argon2Config ...utils.Argon2Config) (bool, error)
}

var _ IPassword = &passwordImpl{}

// Password 实例化工具类方法
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

func (*passwordImpl) RSAEncrypt(plaintext string, publicKey []byte) (string, error) {
	return utils.RSA().Encrypt(plaintext, publicKey)
}

func (*passwordImpl) RSADecrypt(ciphertext string, privateKey []byte) (string, error) {
	return utils.RSA().Decrypt(ciphertext, privateKey)
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
