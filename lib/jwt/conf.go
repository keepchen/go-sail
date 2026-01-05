package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"strings"

	"github.com/keepchen/go-sail/v3/constants"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

// Conf 配置信息
//
// <yaml example>
//
// jwt_conf:
//
//	 enable: false
//		public_key:
//		private_key:
//		algorithm: RS256
//		hmac_secret: example
//		token_issuer: authority
//
// <toml example>
//
// # ::jwt密钥配置::
//
// [jwt_conf]
//
// enable = false
//
// # 公钥文件或字符串
//
// public_key = ""
//
// # 私钥文件或字符串
//
// private_key = ""
//
// # 加密算法: RS256 | RS512 | HS512
//
// algorithm = "RS256"
//
// # 密钥，当algorithm = "HS512"时需要配置此项
//
// hmac_secret = "example"
//
// # 令牌颁发者
//
// token_issuer = "authority"
type Conf struct {
	Enable            bool   `yaml:"enable" toml:"enable" json:"enable"`                   //是否启用
	PublicKey         string `yaml:"public_key" toml:"public_key" json:"public_key"`       //公钥字符串或公钥文件地址(支持rsa、ecdsa、ed25519)
	PrivateKey        string `yaml:"private_key" toml:"private_key" json:"private_key"`    //私钥字符串或私钥文件地址(支持rsa、ecdsa、ed25519)
	Algorithm         string `yaml:"algorithm" toml:"algorithm" json:"algorithm"`          //加密算法: RS256 | RS512 | HS512 | EdDSA | ES256 | ES384 | ES512
	HmacSecret        string `yaml:"hmac_secret" toml:"hmac_secret" json:"hmac_secret"`    //密钥
	TokenIssuer       string `yaml:"token_issuer" toml:"token_issuer" json:"token_issuer"` //令牌颁发者
	rsaPrivateKey     *rsa.PrivateKey
	rsaPublicKey      *rsa.PublicKey
	ed25519PrivateKey crypto.PrivateKey
	ed25519PublicKey  crypto.PublicKey
	ecdsaPrivateKey   *ecdsa.PrivateKey
	ecdsaPublicKey    *ecdsa.PublicKey
}

// MustLoad 载入配置
//
// 载入并解析公私钥，公私钥必须都存在，否则不会解析
func (c *Conf) MustLoad() {
	if len(c.PrivateKey) == 0 || len(c.PublicKey) == 0 {
		return
	}

	var (
		publicBeginStr  = constants.PublicKeyBeginStr
		publicEndStr    = constants.PublicKeyEndStr
		privateBeginStr = constants.PrivateKeyBeginStr
		privateEndStr   = constants.PrivateKeyEndStr
	)

	//当字符串中存在标准前缀头，说明是标准公钥字符，不做任何处理
	if !strings.HasPrefix(c.PublicKey, publicBeginStr) {
		//如果文件存在，从文件读取
		if fileExists(c.PublicKey) {
			contents, err := fileGetContents(c.PublicKey)
			if err == nil {
				c.PublicKey = string(contents)
			}
		}
	}

	//清理前后缀
	c.PublicKey = strings.Replace(c.PublicKey, publicBeginStr, "", 1)
	c.PublicKey = strings.Replace(c.PublicKey, publicEndStr, "", 1)
	c.PublicKey = publicBeginStr + "\n" + wordwrap(c.PublicKey, 64, "\n") + "\n" + publicEndStr

	//当字符串中存在标准前缀头，说明是标准私钥字符，不做任何处理
	if !strings.HasPrefix(c.PrivateKey, privateBeginStr) {
		//如果文件存在，从文件读取
		if fileExists(c.PrivateKey) {
			contents, err := fileGetContents(c.PrivateKey)
			if err == nil {
				c.PrivateKey = string(contents)
			}
		}
	}

	//清理前后缀
	c.PrivateKey = strings.Replace(c.PrivateKey, privateBeginStr, "", 1)
	c.PrivateKey = strings.Replace(c.PrivateKey, privateEndStr, "", 1)
	c.PrivateKey = privateBeginStr + "\n" + wordwrap(c.PrivateKey, 64, "\n") + "\n" + privateEndStr

	if len(c.PrivateKey) == 0 || len(c.PublicKey) == 0 {
		return
	}

	switch c.Algorithm {
	case SigningMethodRS256.String(), SigningMethodRS512.String():
		pri, err := jwtLib.ParseRSAPrivateKeyFromPEM([]byte(c.PrivateKey))
		if err != nil {
			panic(err)
		}
		pub, err := jwtLib.ParseRSAPublicKeyFromPEM([]byte(c.PublicKey))
		if err != nil {
			panic(err)
		}
		c.rsaPrivateKey = pri
		c.rsaPublicKey = pub
	case SigningMethodEdDSA.String():
		pri, err := jwtLib.ParseEdPrivateKeyFromPEM([]byte(c.PrivateKey))
		if err != nil {
			panic(err)
		}
		pub, err := jwtLib.ParseEdPublicKeyFromPEM([]byte(c.PublicKey))
		if err != nil {
			panic(err)
		}
		c.ed25519PrivateKey = pri
		c.ed25519PublicKey = pub
	case SigningMethodES256.String(), SigningMethodES384.String(), SigningMethodES512.String():
		pri, err := jwtLib.ParseECPrivateKeyFromPEM([]byte(c.PrivateKey))
		if err != nil {
			panic(err)
		}
		pub, err := jwtLib.ParseECPublicKeyFromPEM([]byte(c.PublicKey))
		if err != nil {
			panic(err)
		}
		c.ecdsaPrivateKey = pri
		c.ecdsaPublicKey = pub
	}
}

// Load 载入配置
//
// 载入并解析公私钥，公私钥只要任意一个存在就会解析
func (c *Conf) Load() {
	var (
		publicBeginStr  = constants.PublicKeyBeginStr
		publicEndStr    = constants.PublicKeyEndStr
		privateBeginStr = constants.PrivateKeyBeginStr
		privateEndStr   = constants.PrivateKeyEndStr
	)

	if len(c.PublicKey) != 0 {
		//当字符串中存在标准前缀头，说明是标准公钥字符，不做任何处理
		if !strings.HasPrefix(c.PublicKey, publicBeginStr) {
			//如果文件存在，从文件读取
			if fileExists(c.PublicKey) {
				contents, err := fileGetContents(c.PublicKey)
				if err == nil {
					c.PublicKey = string(contents)
				}
			}
		}

		//清理前后缀
		c.PublicKey = strings.Replace(c.PublicKey, publicBeginStr, "", 1)
		c.PublicKey = strings.Replace(c.PublicKey, publicEndStr, "", 1)
		c.PublicKey = publicBeginStr + "\n" + wordwrap(c.PublicKey, 64, "\n") + "\n" + publicEndStr

		switch c.Algorithm {
		case SigningMethodRS256.String(), SigningMethodRS512.String():
			pub, err := jwtLib.ParseRSAPublicKeyFromPEM([]byte(c.PublicKey))
			if err != nil {
				panic(err)
			}
			c.rsaPublicKey = pub
		case SigningMethodEdDSA.String():
			pub, err := jwtLib.ParseEdPublicKeyFromPEM([]byte(c.PublicKey))
			if err != nil {
				panic(err)
			}
			c.ed25519PublicKey = pub
		case SigningMethodES256.String(), SigningMethodES384.String(), SigningMethodES512.String():
			pub, err := jwtLib.ParseECPublicKeyFromPEM([]byte(c.PublicKey))
			if err != nil {
				panic(err)
			}
			c.ecdsaPublicKey = pub
		}
	}

	if len(c.PrivateKey) != 0 {
		//当字符串中存在标准前缀头，说明是标准私钥字符，不做任何处理
		if !strings.HasPrefix(c.PrivateKey, privateBeginStr) {
			//如果文件存在，从文件读取
			if fileExists(c.PrivateKey) {
				contents, err := fileGetContents(c.PrivateKey)
				if err == nil {
					c.PrivateKey = string(contents)
				}
			}
		}

		//清理前后缀
		c.PrivateKey = strings.Replace(c.PrivateKey, privateBeginStr, "", 1)
		c.PrivateKey = strings.Replace(c.PrivateKey, privateEndStr, "", 1)
		c.PrivateKey = privateBeginStr + "\n" + wordwrap(c.PrivateKey, 64, "\n") + "\n" + privateEndStr

		switch c.Algorithm {
		case SigningMethodRS256.String(), SigningMethodRS512.String():
			pri, err := jwtLib.ParseRSAPrivateKeyFromPEM([]byte(c.PrivateKey))
			if err != nil {
				panic(err)
			}
			c.rsaPrivateKey = pri
		case SigningMethodEdDSA.String():
			pri, err := jwtLib.ParseEdPrivateKeyFromPEM([]byte(c.PrivateKey))
			if err != nil {
				panic(err)
			}
			c.ed25519PrivateKey = pri
		case SigningMethodES256.String(), SigningMethodES384.String(), SigningMethodES512.String():
			pri, err := jwtLib.ParseECPrivateKeyFromPEM([]byte(c.PrivateKey))
			if err != nil {
				panic(err)
			}
			c.ecdsaPrivateKey = pri
		}
	}
}

// GetRSAPrivateKeyObj 获取rsa私钥对象
func (c *Conf) GetRSAPrivateKeyObj() *rsa.PrivateKey {
	return c.rsaPrivateKey
}

// GetRSAPublicKeyObj 获取rsa公钥对象
func (c *Conf) GetRSAPublicKeyObj() *rsa.PublicKey {
	return c.rsaPublicKey
}

// GetED25519PrivateKeyObj 获取ed25519私钥对象
func (c *Conf) GetED25519PrivateKeyObj() crypto.PrivateKey {
	return c.ed25519PrivateKey
}

// GetED25519PublicKeyObj 获取ed25519公钥对象
func (c *Conf) GetED25519PublicKeyObj() crypto.PublicKey {
	return c.ed25519PublicKey
}

// GetECDSAPrivateKeyObj 获取ecdsa私钥对象
func (c *Conf) GetECDSAPrivateKeyObj() *ecdsa.PrivateKey {
	return c.ecdsaPrivateKey
}

// GetECDSAPublicKeyObj 获取ecdsa公钥对象
func (c *Conf) GetECDSAPublicKeyObj() *ecdsa.PublicKey {
	return c.ecdsaPublicKey
}
