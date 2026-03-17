package utils

import "golang.org/x/crypto/bcrypt"

type bcryptImpl struct{}

// IBCrypt bcrypt接口
type IBCrypt interface {
	// HashMake 生成哈希字符
	//
	// plaintext - 明文字符
	//
	// costs - 权重，范围∈[4,31]，默认为10
	HashMake(plaintext string, costs ...int) (string, error)
	// HashCheck 验证哈希字符
	//
	// plaintext - 明文字符
	//
	// hashed - 已哈希的字符
	HashCheck(plaintext, hashed string) (bool, error)
}

var _ IBCrypt = &bcryptImpl{}

// Bcrypt 实例化bcrypt哈希
func Bcrypt() IBCrypt {
	return &bcryptImpl{}
}

// HashMake 生成哈希字符
//
// plaintext - 明文字符
//
// costs - 权重，范围∈[4,31]，默认为10
func (bcryptImpl) HashMake(plaintext string, costs ...int) (string, error) {
	var cost = bcrypt.DefaultCost
	if len(costs) > 0 {
		cost = costs[0]
	}
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), cost)

	return string(bytes), err
}

// HashCheck 验证哈希字符
//
// plaintext - 明文字符
//
// hashed - 已哈希的字符
func (bcryptImpl) HashCheck(plaintext, hashed string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plaintext))

	return err == nil, err
}
