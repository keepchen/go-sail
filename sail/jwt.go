package sail

import (
	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/sail/config"
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
	ValidToken(token string) (bool, error)
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
func (jwtImpl) ValidToken(token string) (bool, error) {
	conf := config.Get()
	_, err := jwt.VerifyFromMap(token, *conf.JwtConf)
	if err != nil {
		return false, err
	}

	return err == nil, err
}
