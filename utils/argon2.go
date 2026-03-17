package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"math/rand/v2"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Config 配置参数
type Argon2Config struct {
	Memory      uint32 // 内存使用量（KB）
	Iterations  uint32 // 迭代次数
	Parallelism uint8  // 并行度
	SaltLength  uint32 // 盐长度
	KeyLength   uint32 // 输出密钥长度
}

type argon2Impl struct {
	config Argon2Config
}

// IArgon2 argon2接口
type IArgon2 interface {
	// HashMake 生成哈希字符
	//
	// plaintext - 明文字符
	HashMake(plaintext string) (string, error)
	// HashCheck 验证哈希字符
	//
	// plaintext - 明文字符
	//
	// hashed - 已哈希的字符
	HashCheck(plaintext, hashed string) (bool, error)
}

var _ IArgon2 = &argon2Impl{}

// 默认的Argon2配置
var defaultArgon2Config = Argon2Config{
	Memory:      64 * 1024, //64MB
	Iterations:  3,
	Parallelism: 4,
	SaltLength:  16,
	KeyLength:   32,
}

const (
	argon2EncodeFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	argon2DecodeFormat = "m=%d,t=%d,p=%d"
)

// Argon2 实例化argon2算法
func Argon2(config ...Argon2Config) IArgon2 {
	if len(config) > 0 {
		return &argon2Impl{config[0]}
	}
	return &argon2Impl{
		config: defaultArgon2Config,
	}
}

// HashMake 生成哈希字符
//
// plaintext - 明文字符
func (a *argon2Impl) HashMake(plaintext string) (string, error) {
	salt := make([]byte, a.config.SaltLength)
	if _, err := (&rand.ChaCha8{}).Read(salt); err != nil {
		return "", err
	}
	hashed := argon2.IDKey([]byte(plaintext), salt,
		a.config.Iterations,
		a.config.Memory,
		a.config.Parallelism,
		a.config.KeyLength)

	b64salt := base64.RawStdEncoding.EncodeToString(salt)
	b64hash := base64.RawStdEncoding.EncodeToString(hashed)

	result := fmt.Sprintf(argon2EncodeFormat, argon2.Version, a.config.Memory, a.config.Iterations, a.config.Parallelism, b64salt, b64hash)

	return result, nil
}

// HashCheck 验证哈希字符
//
// plaintext - 明文字符
//
// hashed - 已哈希的字符
func (a *argon2Impl) HashCheck(plaintext, hashed string) (bool, error) {
	parts := strings.Split(hashed, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hashed format: %s", hashed)
	}
	var version int
	var cfg Argon2Config
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	_, err = fmt.Sscanf(parts[3], argon2DecodeFormat, &cfg.Memory, &cfg.Iterations, &cfg.Parallelism)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	cfg.SaltLength = uint32(len(salt))
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	cfg.KeyLength = uint32(len(hash))
	//重新计算哈希
	computedHash := argon2.IDKey([]byte(plaintext), salt,
		cfg.Iterations, cfg.Memory, cfg.Parallelism, cfg.KeyLength)

	// 使用constant-time比较防止时序攻击
	return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
