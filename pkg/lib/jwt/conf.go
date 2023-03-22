package jwt

import (
	"crypto/rsa"
	"strings"

	"github.com/keepchen/go-sail/pkg/constants"

	jwtLib "github.com/golang-jwt/jwt"
)

// Conf 配置信息
type Conf struct {
	PublicKey   string `yaml:"public_key" toml:"public_key" json:"public_key"`       //公钥字符串或公钥文件地址
	PrivateKey  string `yaml:"private_key" toml:"private_key" json:"private_key"`    //私钥字符串或私钥文件地址
	Algorithm   string `yaml:"algorithm" toml:"algorithm" json:"algorithm"`          //加密算法: RS256 | RS512 | HS512
	HmacSecret  string `yaml:"hmac_secret" toml:"hmac_secret" json:"hmac_secret"`    //密钥
	TokenIssuer string `yaml:"token_issuer" toml:"token_issuer" json:"token_issuer"` //令牌颁发者
	privateKey  *rsa.PrivateKey
	publicKey   *rsa.PublicKey
}

// Load 载入配置
func (c *Conf) Load() {
	if len(c.PrivateKey) == 0 || len(c.PublicKey) == 0 {
		return
	}

	//当字符串中存在标准前缀头，说明是标准公钥字符，不做任何处理
	if !strings.HasPrefix(c.PublicKey, constants.PublicKeyBeginStr) {
		//如果文件存在，从文件读取
		if fileExists(c.PublicKey) {
			contents, err := fileGetContents(c.PublicKey)
			if err != nil {
				c.PublicKey = string(contents)
			}
		} else {
			//如果文件不存在，则转换字符串
			c.PublicKey = constants.PublicKeyBeginStr + "\n" + wordwrap(c.PublicKey, 64, "\n") + "\n" + constants.PublicKeyEndStr
		}
	}

	//当字符串中存在标准前缀头，说明是标准私钥字符，不做任何处理
	if !strings.HasPrefix(c.PrivateKey, constants.PrivateKeyBeginStr) {
		//如果文件存在，从文件读取
		if fileExists(c.PrivateKey) {
			contents, err := fileGetContents(c.PrivateKey)
			if err != nil {
				c.PrivateKey = string(contents)
			}
		} else {
			//如果文件不存在，则转换字符串
			c.PrivateKey = constants.PrivateKeyBeginStr + "\n" + wordwrap(c.PrivateKey, 64, "\n") + "\n" + constants.PrivateKeyEndStr
		}
	}

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

// GetPrivateKeyObj 获取私钥对象
func (c *Conf) GetPrivateKeyObj() *rsa.PrivateKey {
	return c.privateKey
}

// GetPublicKeyObj 获取公钥对象
func (c *Conf) GetPublicKeyObj() *rsa.PublicKey {
	return c.publicKey
}
