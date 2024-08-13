package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	jwtLib "github.com/golang-jwt/jwt"
)

// AppClaims App票据声明
type AppClaims struct {
	Scopes   []string
	ScopeIDs []int64
	Name     string
	Email    string
	//more...
	jwtLib.StandardClaims
}

type MapClaims map[string]interface{}

// Valid 验证token有效性
//
// 此验证方法继承 jwtLib.StandardClaims 的 Valid 验证方法
func (c *MapClaims) Valid() error {
	if c == nil {
		return fmt.Errorf("MapClaims is nil")
	}

	var standardClaim = &jwtLib.StandardClaims{}
	if value, ok := (*c)["aud"]; ok {
		if aud, ok := value.(string); ok {
			standardClaim.Audience = aud
		}
	}
	if value, ok := (*c)["exp"]; ok {
		valueStr := fmt.Sprintf("%v", value)
		if exp, err := strconv.ParseFloat(valueStr, 64); err == nil {
			standardClaim.ExpiresAt = int64(exp)
		}
	}
	if value, ok := (*c)["jti"]; ok {
		if jti, ok := value.(string); ok {
			standardClaim.Id = jti
		}
	}
	if value, ok := (*c)["iat"]; ok {
		valueStr := fmt.Sprintf("%v", value)
		if iat, err := strconv.ParseFloat(valueStr, 64); err == nil {
			standardClaim.IssuedAt = int64(iat)
		}
	}
	if value, ok := (*c)["iss"]; ok {
		if iss, ok := value.(string); ok {
			standardClaim.Issuer = iss
		}
	}
	if value, ok := (*c)["nbf"]; ok {
		valueStr := fmt.Sprintf("%v", value)
		if nbf, err := strconv.ParseFloat(valueStr, 64); err == nil {
			standardClaim.NotBefore = int64(nbf)
		}
	}
	if value, ok := (*c)["sub"]; ok {
		if sub, ok := value.(string); ok {
			standardClaim.Subject = sub
		}
	}

	return standardClaim.Valid()
}

// MergeStandardClaims
//
// 合并标准字段
//
// 如果传入的自定义字段在标准字段中存在，则用自定义字段覆盖标准字段
func MergeStandardClaims(fields map[string]interface{}) MapClaims {
	now := time.Now()
	defaultClaims := MapClaims{
		"jti": uuid.New().String(),
		"iat": now.Unix(),
		"exp": now.Add(time.Hour * 24).Unix(),
		"nbf": now.Unix(),
		"iss": "Go-Sail",
	}

	//override
	for k, v := range fields {
		defaultClaims[k] = v
	}

	return defaultClaims
}
