package jwt

import (
	"errors"
	"fmt"

	jwtLib "github.com/golang-jwt/jwt"
)

const (
	defaultTokenIssuer = "authority"
)

// SigningMethod 签名方法
type SigningMethod string

const (
	SigningMethodRS256 SigningMethod = "RS256" //SigningMethodRS256 rsa256方法
	SigningMethodRS512 SigningMethod = "RS512" //SigningMethodRS512 rsa512方法
	SigningMethodHS512 SigningMethod = "HS512" //SigningMethodHS512 hmac方法
)

func (sm SigningMethod) getSigningMethod() jwtLib.SigningMethod {
	return jwtLib.GetSigningMethod(string(sm))
}

// Sign 签名
func Sign(claim jwtLib.Claims, conf Conf) (string, error) {
	appClaims, ok := claim.(AppClaims)
	if !ok {
		return "", errors.New("not supported claims")
	}

	switch conf.Algorithm {
	case string(SigningMethodRS256):
		return appClaims.GetToken(SigningMethodRS256, conf.privateKey)
	case string(SigningMethodRS512):
		return appClaims.GetToken(SigningMethodRS512, conf.privateKey)
	case string(SigningMethodHS512):
		return appClaims.GetToken(SigningMethodHS512, []byte(conf.HmacSecret))
	default:
		return "", errors.New("jwt secret not config")
	}
}

// SignWithMap 签名
func SignWithMap(claims MapClaims, conf Conf) (string, error) {
	switch conf.Algorithm {
	case string(SigningMethodRS256):
		return claims.GetToken(SigningMethodRS256, conf.privateKey)
	case string(SigningMethodRS512):
		return claims.GetToken(SigningMethodRS512, conf.privateKey)
	case string(SigningMethodHS512):
		return claims.GetToken(SigningMethodHS512, []byte(conf.HmacSecret))
	default:
		return "", errors.New("jwt secret not config")
	}
}

// Verify 验证
func Verify(tokenString string, conf Conf) (AppClaims, error) {
	token, err := jwtLib.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwtLib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodHMAC); ok {
			return []byte(conf.HmacSecret), nil
		}

		if _, ok := token.Method.(*jwtLib.SigningMethodRSA); ok {
			//私钥加密，公钥解密，因此这里返回公钥
			return conf.privateKey.Public(), nil
		}

		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})

	if err != nil {
		return AppClaims{}, err
	}

	if claims, ok := token.Claims.(*AppClaims); ok && token.Valid {
		return *claims, nil
	}

	return AppClaims{}, errors.New("token verify failed")
}

// VerifyFromMap 验证
func VerifyFromMap(tokenString string, conf Conf) (MapClaims, error) {
	token, err := jwtLib.ParseWithClaims(tokenString, &MapClaims{}, func(token *jwtLib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtLib.SigningMethodHMAC); ok {
			return []byte(conf.HmacSecret), nil
		}

		if _, ok := token.Method.(*jwtLib.SigningMethodRSA); ok {
			//私钥加密，公钥解密，因此这里返回公钥
			return conf.privateKey.Public(), nil
		}

		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	})

	if err != nil {
		return MapClaims{}, err
	}

	if claims, ok := token.Claims.(*MapClaims); ok && token.Valid {
		return *claims, nil
	}

	return MapClaims{}, errors.New("token verify failed")
}

// GetToken 获取token
func (c *AppClaims) GetToken(method SigningMethod, secret interface{}) (string, error) {
	if secret == nil {
		return "", errors.New("unsupported secret")
	}
	sm := method.getSigningMethod()
	token := jwtLib.NewWithClaims(sm, *c)
	tokenStr, err := token.SignedString(secret)

	return tokenStr, err
}

// GetToken 获取token
func (c *MapClaims) GetToken(method SigningMethod, secret interface{}) (string, error) {
	if secret == nil {
		return "", errors.New("unsupported secret")
	}
	sm := method.getSigningMethod()
	token := jwtLib.NewWithClaims(sm, c)
	tokenStr, err := token.SignedString(secret)

	return tokenStr, err
}
