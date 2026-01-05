package jwt

import (
	"errors"
	"fmt"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

const (
	defaultTokenIssuer = "Go-Sail"
)

// SigningMethod 签名方法
type SigningMethod string

const (
	SigningMethodRS256 SigningMethod = "RS256" //SigningMethodRS256 rsa256方法
	SigningMethodRS512 SigningMethod = "RS512" //SigningMethodRS512 rsa512方法
	SigningMethodHS512 SigningMethod = "HS512" //SigningMethodHS512 hmac方法
	SigningMethodEdDSA SigningMethod = "EdDSA" //SigningMethodEdDSA ed25519方法
	SigningMethodES256 SigningMethod = "ES256" //SigningMethodES256 ecdsa256方法
	SigningMethodES384 SigningMethod = "ES384" //SigningMethodES384 ecdsa384方法
	SigningMethodES512 SigningMethod = "ES512" //SigningMethodES521 ecdsa512方法
)

func (sm SigningMethod) String() string {
	return string(sm)
}

func (sm SigningMethod) getSigningMethod() jwtLib.SigningMethod {
	switch sm {
	case SigningMethodRS256:
		return jwtLib.SigningMethodRS256
	case SigningMethodRS512:
		return jwtLib.SigningMethodRS512
	case SigningMethodHS512:
		return jwtLib.SigningMethodHS512
	case SigningMethodEdDSA:
		return jwtLib.SigningMethodEdDSA
	case SigningMethodES256:
		return jwtLib.SigningMethodES256
	case SigningMethodES384:
		return jwtLib.SigningMethodES384
	case SigningMethodES512:
		return jwtLib.SigningMethodES512
	default:
		return jwtLib.GetSigningMethod(sm.String())
	}
}

// Sign 签名
func Sign(claim jwtLib.Claims, conf Conf) (string, error) {
	appClaims, ok := claim.(AppClaims)
	if !ok {
		return "", errors.New("not supported claims")
	}

	if len(appClaims.Issuer) == 0 {
		appClaims.Issuer = defaultTokenIssuer
	}

	switch conf.Algorithm {
	case string(SigningMethodRS256):
		return appClaims.GetToken(SigningMethodRS256, conf.rsaPrivateKey)
	case string(SigningMethodRS512):
		return appClaims.GetToken(SigningMethodRS512, conf.rsaPrivateKey)
	case string(SigningMethodHS512):
		return appClaims.GetToken(SigningMethodHS512, []byte(conf.HmacSecret))
	case string(SigningMethodEdDSA):
		return appClaims.GetToken(SigningMethodEdDSA, conf.ed25519PrivateKey)
	case string(SigningMethodES256):
		return appClaims.GetToken(SigningMethodES256, conf.ecdsaPrivateKey)
	case string(SigningMethodES384):
		return appClaims.GetToken(SigningMethodES384, conf.ecdsaPrivateKey)
	case string(SigningMethodES512):
		return appClaims.GetToken(SigningMethodES512, conf.ecdsaPrivateKey)
	default:
		return "", errors.New("algorithm is not set, it must be one of RS256, RS512, HS512")
	}
}

// SignWithMap 签名
func SignWithMap(claims MapClaims, conf Conf) (string, error) {
	if val, ok := claims["iss"]; !ok || len(fmt.Sprintf("%v", val)) == 0 {
		if len(conf.TokenIssuer) > 0 {
			claims["iss"] = conf.TokenIssuer
		} else {
			claims["iss"] = defaultTokenIssuer
		}
	}

	switch conf.Algorithm {
	case string(SigningMethodRS256):
		return claims.GetToken(SigningMethodRS256, conf.rsaPrivateKey)
	case string(SigningMethodRS512):
		return claims.GetToken(SigningMethodRS512, conf.rsaPrivateKey)
	case string(SigningMethodHS512):
		return claims.GetToken(SigningMethodHS512, []byte(conf.HmacSecret))
	case string(SigningMethodEdDSA):
		return claims.GetToken(SigningMethodEdDSA, conf.ed25519PrivateKey)
	case string(SigningMethodES256):
		return claims.GetToken(SigningMethodES256, conf.ecdsaPrivateKey)
	case string(SigningMethodES384):
		return claims.GetToken(SigningMethodES384, conf.ecdsaPrivateKey)
	case string(SigningMethodES512):
		return claims.GetToken(SigningMethodES512, conf.ecdsaPrivateKey)
	default:
		return "", errors.New("algorithm is not set, it must be one of RS256, RS512, HS512")
	}
}

// Verify 验证
func Verify(tokenString string, conf Conf) (AppClaims, error) {
	token, err := jwtLib.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwtLib.Token) (any, error) {
		if token.Method.Alg() != conf.Algorithm {
			return nil, fmt.Errorf(
				"unexpected signing algorithm: %s",
				token.Method.Alg(),
			)
		}

		switch token.Method.(type) {
		case *jwtLib.SigningMethodHMAC:
			return []byte(conf.HmacSecret), nil
		case *jwtLib.SigningMethodRSA:
			return conf.rsaPublicKey, nil
		case *jwtLib.SigningMethodECDSA:
			return conf.ecdsaPublicKey, nil
		case *jwtLib.SigningMethodEd25519:
			return conf.ed25519PublicKey, nil
		default:
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
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
	token, err := jwtLib.ParseWithClaims(tokenString, &MapClaims{}, func(token *jwtLib.Token) (any, error) {
		if token.Method.Alg() != conf.Algorithm {
			return nil, fmt.Errorf(
				"unexpected signing algorithm: %s",
				token.Method.Alg(),
			)
		}

		switch token.Method.(type) {
		case *jwtLib.SigningMethodHMAC:
			return []byte(conf.HmacSecret), nil
		case *jwtLib.SigningMethodRSA:
			return conf.rsaPublicKey, nil
		case *jwtLib.SigningMethodECDSA:
			return conf.ecdsaPublicKey, nil
		case *jwtLib.SigningMethodEd25519:
			return conf.ed25519PublicKey, nil
		default:
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
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
func (c *AppClaims) GetToken(method SigningMethod, secret any) (string, error) {
	if secret == nil {
		return "", errors.New("secret is nil")
	}
	sm := method.getSigningMethod()
	token := jwtLib.NewWithClaims(sm, *c)
	tokenStr, err := token.SignedString(secret)

	return tokenStr, err
}

// GetToken 获取token
func (c *MapClaims) GetToken(method SigningMethod, secret any) (string, error) {
	if secret == nil {
		return "", errors.New("secret is nil")
	}
	sm := method.getSigningMethod()
	token := jwtLib.NewWithClaims(sm, c)
	tokenStr, err := token.SignedString(secret)

	return tokenStr, err
}
