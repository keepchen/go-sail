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
// issuer 颁发者
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
	for _, otherField := range otherFields {
		for k, v := range otherField {
			baseMap[k] = v
		}
	}
	mp := jwt.MergeStandardClaims(baseMap)

	return jwt.SignWithMap(mp, *conf.JwtConf)
}

// ValidToken 验证JWT令牌
func (jwtImpl) ValidToken(token string) (bool, error) {
	conf := config.Get()
	mp, err := jwt.VerifyFromMap(token, *conf.JwtConf)
	if err != nil {
		return false, err
	}
	err = mp.Valid()

	return err == nil, err
}
