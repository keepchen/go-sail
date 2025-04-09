package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

// AppClaims App票据声明
type AppClaims struct {
	Scopes   []string
	ScopeIDs []int64
	Name     string
	Email    string
	//more...
	jwtLib.RegisteredClaims
}

type MapClaims map[string]interface{}

func (c *MapClaims) GetExpirationTime() (*jwtLib.NumericDate, error) {
	if value, ok := (*c)["exp"]; ok {
		valueStr := fmt.Sprintf("%v", value)
		if exp, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return jwtLib.NewNumericDate(time.Unix(int64(exp), 0)), nil
		}
	}

	return nil, errors.New("exp not registered")
}

func (c *MapClaims) GetIssuedAt() (*jwtLib.NumericDate, error) {
	if value, ok := (*c)["iat"]; ok {
		valueStr := fmt.Sprintf("%v", value)
		if iat, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return jwtLib.NewNumericDate(time.Unix(int64(iat), 0)), nil
		}
	}

	return nil, errors.New("iat not registered")
}

func (c *MapClaims) GetNotBefore() (*jwtLib.NumericDate, error) {
	if value, ok := (*c)["nbf"]; ok {
		valueStr := fmt.Sprintf("%v", value)
		if nbf, err := strconv.ParseFloat(valueStr, 64); err == nil {
			return jwtLib.NewNumericDate(time.Unix(int64(nbf), 0)), nil
		}
	}

	return nil, errors.New("nbf not registered")
}

func (c *MapClaims) GetIssuer() (string, error) {
	if value, ok := (*c)["iss"]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "", errors.New("iss not registered")
}

func (c *MapClaims) GetSubject() (string, error) {
	if value, ok := (*c)["sub"]; ok {
		return fmt.Sprintf("%v", value), nil
	}

	return "", errors.New("sub not registered")
}

func (c *MapClaims) GetAudience() (jwtLib.ClaimStrings, error) {
	if value, ok := (*c)["aud"]; ok {
		return jwtLib.ClaimStrings{fmt.Sprintf("%v", value)}, nil
	}

	return jwtLib.ClaimStrings{""}, errors.New("aud not registered")
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
