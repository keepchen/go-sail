package sail

import "github.com/keepchen/go-sail/v3/utils"

type passwordImpl struct{}

// IPassword 密码工具类接口
type IPassword interface {
	// AesGCMEncrypt aes加密
	//
	// 使用GCM
	//
	// key应该是一个16或24或32位长度的字符
	AesGCMEncrypt(plaintext, key string) (string, error)
	// AesGCMDecrypt aes解密
	//
	// 使用GCM
	//
	// key应该是一个16或24或32位长度的字符
	AesGCMDecrypt(ciphertext, key string) (string, error)
	// SM4GCMEncrypt GCM加密
	//
	// hexKey 16进制key 长度32位
	//
	// raw 待加密内容
	SM4GCMEncrypt(plaintext, key string) (string, error)
	// SM4GCMDecrypt GCM解密
	//
	// hexKey 16进制key 长度32位
	//
	// base64Raw 加密内容 base64格式
	SM4GCMDecrypt(ciphertext, key string) (string, error)
	// RSAEncrypt rsa加密
	//
	// publicKey格式支持：
	//
	// - 标准格式
	//
	// - 不带前后缀的一行字符串
	RSAEncrypt(plaintext string, publicKey []byte) (string, error)
	// RSADecrypt rsa解密
	//
	// privateKey格式支持：
	//
	// - 标准格式
	//
	// - 不带前后缀的一行字符串
	RSADecrypt(ciphertext string, privateKey []byte) (string, error)
	// BcryptHashMake 生成哈希字符
	//
	// plaintext - 明文字符
	//
	// costs - 权重，范围∈[4,31]，默认为10
	BcryptHashMake(plaintext string, costs ...int) (string, error)
	// BcryptHashCheck 验证哈希字符
	//
	// plaintext - 明文字符
	//
	// hashed - 已哈希的字符
	BcryptHashCheck(plaintext string, hash string) (bool, error)
	// Argon2HashMake 生成哈希字符
	//
	// plaintext - 明文字符
	//
	// argon2Config - 配置参数
	Argon2HashMake(plaintext string, argon2Config ...utils.Argon2Config) (string, error)
	// Argon2HashCheck 验证哈希字符
	//
	// plaintext - 明文字符
	//
	// hashed - 已哈希的字符
	//
	// argon2Config - 配置参数
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
