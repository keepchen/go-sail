package sail

import (
	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/keepchen/go-sail/v3/utils"
)

// IJwt Jwt接口
type IJwt interface {
	// MakeToken 生成JWT令牌
	//
	// uid 用户id
	//
	// exp 令牌有效期，秒级时间戳
	//
	// otherFields 其他可扩展字段
	MakeToken(uid string, exp int64, otherFields ...map[string]interface{}) (string, error)
	// ValidToken 验证JWT令牌
	//
	// # ⚠️注意⚠️
	//
	// 此方法继承 jwt.Validator 的 Validate 方法，因此只验证了存在且格式正确的字段，
	// 如果调用方有其他的验证规则，需要自行处理验证逻辑。
	ValidToken(token string) (bool, jwt.MapClaims, error)
	//Encrypt 加密字符
	//
	// 此方法可以与 Decrypt 方法对应进行加解密
	//
	// # Note
	//
	// 此方法调用 utils.RSA 方法并使用 jwt.Conf 配置中的公钥PublicKey进行加密
	Encrypt(plaintext string) (string, error)
	//Decrypt 解密字符
	//
	// 此方法可以与 Encrypt 方法对应进行加解密
	//
	// # Note
	//
	// 此方法调用 utils.RSA 方法并使用 jwt.Conf 配置中的私钥PrivateKey进行解密
	Decrypt(encryptedStr string) (string, error)
}

type jwtImpl struct{}

var _ IJwt = &jwtImpl{}

// JWT jwt实例
func JWT() IJwt {
	return &jwtImpl{}
}

// MakeToken 生成JWT令牌
//
// uid 用户id
//
// exp 令牌有效期，秒级时间戳
//
// otherFields 其他可扩展字段
func (jwtImpl) MakeToken(uid string, exp int64, otherFields ...map[string]interface{}) (string, error) {
	conf := config.Get()
	baseMap := map[string]interface{}{
		"uid": uid,
		"iss": conf.JwtConf.TokenIssuer,
		"exp": exp,
	}
	//可能会存在覆盖的情况
	for _, otherField := range otherFields {
		for k, v := range otherField {
			baseMap[k] = v
		}
	}
	mp := jwt.MergeStandardClaims(baseMap)

	return jwt.SignWithMap(mp, *conf.JwtConf)
}

// ValidToken 验证JWT令牌
//
// # ⚠️注意⚠️
//
// 此方法继承 jwt.Validator 的 Validate 方法，因此只验证了存在且格式正确的字段，
// 如果调用方有其他的验证规则，需要自行处理验证逻辑。
func (jwtImpl) ValidToken(token string) (bool, jwt.MapClaims, error) {
	conf := config.Get()
	mp, err := jwt.VerifyFromMap(token, *conf.JwtConf)
	if err != nil {
		return false, mp, err
	}

	return err == nil, mp, err
}

// Encrypt 加密字符
//
// 此方法可以与 Decrypt 方法对应进行加解密
//
// # Note
//
// 此方法调用 utils.RSA 方法并使用 jwt.Conf 配置中的公钥PublicKey进行加密
func (jwtImpl) Encrypt(plaintext string) (string, error) {
	conf := config.Get()
	return utils.RSA().Encrypt(plaintext, []byte(conf.JwtConf.PublicKey))
}

// Decrypt 解密字符
//
// 此方法可以与 Encrypt 方法对应进行加解密
//
// # Note
//
// 此方法调用 utils.RSA 方法并使用 jwt.Conf 配置中的私钥PrivateKey进行解密
func (jwtImpl) Decrypt(encryptedStr string) (string, error) {
	conf := config.Get()
	return utils.RSA().Decrypt(encryptedStr, []byte(conf.JwtConf.PrivateKey))
}
