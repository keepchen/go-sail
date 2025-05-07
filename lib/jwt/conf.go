package jwt

import (
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
	Enable      bool   `yaml:"enable" toml:"enable" json:"enable"`                   //是否启用
	PublicKey   string `yaml:"public_key" toml:"public_key" json:"public_key"`       //公钥字符串或公钥文件地址
	PrivateKey  string `yaml:"private_key" toml:"private_key" json:"private_key"`    //私钥字符串或私钥文件地址
	Algorithm   string `yaml:"algorithm" toml:"algorithm" json:"algorithm"`          //加密算法: RS256 | RS512 | HS512
	HmacSecret  string `yaml:"hmac_secret" toml:"hmac_secret" json:"hmac_secret"`    //密钥
	TokenIssuer string `yaml:"token_issuer" toml:"token_issuer" json:"token_issuer"` //令牌颁发者
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

// MustLoad 载入配置
//
// 载入并解析公私钥，公私钥必须都存在，否则不会解析
func (c *Conf) MustLoad() {
	if len(c.PrivateKey) == 0 || len(c.PublicKey) == 0 {
		return
	}

	//当字符串中存在标准前缀头，说明是标准公钥字符，不做任何处理
	if !strings.HasPrefix(c.PublicKey, constants.PublicKeyBeginStr) {
		//如果文件存在，从文件读取
		if fileExists(c.PublicKey) {
			contents, err := fileGetContents(c.PublicKey)
			if err == nil {
				c.PublicKey = string(contents)
			}
		}
	}

	//清理前后缀
	c.PublicKey = strings.Replace(c.PublicKey, constants.PublicKeyBeginStr, "", 1)
	c.PublicKey = strings.Replace(c.PublicKey, constants.PublicKeyEndStr, "", 1)
	c.PublicKey = constants.PublicKeyBeginStr + "\n" + wordwrap(c.PublicKey, 64, "\n") + "\n" + constants.PublicKeyEndStr

	//当字符串中存在标准前缀头，说明是标准私钥字符，不做任何处理
	if !strings.HasPrefix(c.PrivateKey, constants.PrivateKeyBeginStr) {
		//如果文件存在，从文件读取
		if fileExists(c.PrivateKey) {
			contents, err := fileGetContents(c.PrivateKey)
			if err == nil {
				c.PrivateKey = string(contents)
			}
		}
	}

	//清理前后缀
	c.PrivateKey = strings.Replace(c.PrivateKey, constants.PrivateKeyBeginStr, "", 1)
	c.PrivateKey = strings.Replace(c.PrivateKey, constants.PrivateKeyEndStr, "", 1)
	c.PrivateKey = constants.PrivateKeyBeginStr + "\n" + wordwrap(c.PrivateKey, 64, "\n") + "\n" + constants.PrivateKeyEndStr

	if len(c.PrivateKey) == 0 || len(c.PublicKey) == 0 {
		return
	}

	pri, err := jwtLib.ParseRSAPrivateKeyFromPEM([]byte(c.PrivateKey))
	if err != nil {
		panic(err)
	}
	pub, err := jwtLib.ParseRSAPublicKeyFromPEM([]byte(c.PublicKey))
	if err != nil {
		panic(err)
	}
	c.privateKey = pri
	c.publicKey = pub
}

// Load 载入配置
//
// 载入并解析公私钥，公私钥只要任意一个存在就会解析
func (c *Conf) Load() {
	if len(c.PublicKey) != 0 {
		//当字符串中存在标准前缀头，说明是标准公钥字符，不做任何处理
		if !strings.HasPrefix(c.PublicKey, constants.PublicKeyBeginStr) {
			//如果文件存在，从文件读取
			if fileExists(c.PublicKey) {
				contents, err := fileGetContents(c.PublicKey)
				if err == nil {
					c.PublicKey = string(contents)
				}
			}
		}

		//清理前后缀
		c.PublicKey = strings.Replace(c.PublicKey, constants.PublicKeyBeginStr, "", 1)
		c.PublicKey = strings.Replace(c.PublicKey, constants.PublicKeyEndStr, "", 1)
		c.PublicKey = constants.PublicKeyBeginStr + "\n" + wordwrap(c.PublicKey, 64, "\n") + "\n" + constants.PublicKeyEndStr

		pub, err := jwtLib.ParseRSAPublicKeyFromPEM([]byte(c.PublicKey))
		if err != nil {
			panic(err)
		}
		c.publicKey = pub
	}

	if len(c.PrivateKey) != 0 {
		//当字符串中存在标准前缀头，说明是标准私钥字符，不做任何处理
		if !strings.HasPrefix(c.PrivateKey, constants.PrivateKeyBeginStr) {
			//如果文件存在，从文件读取
			if fileExists(c.PrivateKey) {
				contents, err := fileGetContents(c.PrivateKey)
				if err == nil {
					c.PrivateKey = string(contents)
				}
			}
		}

		//清理前后缀
		c.PrivateKey = strings.Replace(c.PrivateKey, constants.PrivateKeyBeginStr, "", 1)
		c.PrivateKey = strings.Replace(c.PrivateKey, constants.PrivateKeyEndStr, "", 1)
		c.PrivateKey = constants.PrivateKeyBeginStr + "\n" + wordwrap(c.PrivateKey, 64, "\n") + "\n" + constants.PrivateKeyEndStr

		pri, err := jwtLib.ParseRSAPrivateKeyFromPEM([]byte(c.PrivateKey))
		if err != nil {
			panic(err)
		}
		c.privateKey = pri
	}
}

// GetPrivateKeyObj 获取私钥对象
func (c *Conf) GetPrivateKeyObj() *rsa.PrivateKey {
	return c.privateKey
}

// GetPublicKeyObj 获取公钥对象
func (c *Conf) GetPublicKeyObj() *rsa.PublicKey {
	return c.publicKey
}
