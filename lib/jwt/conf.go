package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"strings"

	"github.com/keepchen/go-sail/v3/constants"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

// Conf configuration
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
// # ::JWT key configuration::
//
// [jwt_conf]
//
// enable = false
//
// # Public key file or string
//
// public_key = ""
//
// # Private key file or string
//
// private_key = ""
//
// # Encryption algorithm: RS256 | RS512 | HS512
//
// algorithm = "RS256"
//
// # Secret key, required when algorithm = "HS512"
//
// hmac_secret = "example"
//
// # Token issuer
//
// token_issuer = "authority"
type Conf struct {
	Enable            bool   `yaml:"enable" toml:"enable" json:"enable"`                                     //Whether to enable
	PublicKey         string `yaml:"public_key" toml:"public_key" json:"public_key"`                         //Public key string or file path (supports RSA, ECDSA, and Ed25519)
	PrivateKey        string `yaml:"private_key" toml:"private_key" json:"private_key"`                      //Private key string or file path (supports RSA, ECDSA, and Ed25519)
	Algorithm         string `yaml:"algorithm" toml:"algorithm" json:"algorithm"`                            //Encryption algorithm: RS256 | RS512 | HS512 | EdDSA | ES256 | ES384 | ES512
	HmacSecret        string `yaml:"hmac_secret" toml:"hmac_secret" json:"hmac_secret"`                      //Secret key
	TokenIssuer       string `yaml:"token_issuer" toml:"token_issuer" json:"token_issuer" default:"Go-Sail"` //Token issuer
	rsaPrivateKey     *rsa.PrivateKey
	rsaPublicKey      *rsa.PublicKey
	ed25519PrivateKey crypto.PrivateKey
	ed25519PublicKey  crypto.PublicKey
	ecdsaPrivateKey   *ecdsa.PrivateKey
	ecdsaPublicKey    *ecdsa.PublicKey
}

// MustLoad loads configuration
//
// Load and parse the public and private keys; both must be present for parsing
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

	//A standard begin marker indicates a standard public key string; leave it unchanged
	if !strings.HasPrefix(c.PublicKey, publicBeginStr) {
		//Read from the file if it exists
		if fileExists(c.PublicKey) {
			contents, err := fileGetContents(c.PublicKey)
			if err == nil {
				c.PublicKey = string(contents)
			}
		}
	}

	//Remove the prefix and suffix
	c.PublicKey = strings.Replace(c.PublicKey, publicBeginStr, "", 1)
	c.PublicKey = strings.Replace(c.PublicKey, publicEndStr, "", 1)
	c.PublicKey = publicBeginStr + "\n" + wordwrap(c.PublicKey, 64, "\n") + "\n" + publicEndStr

	//A standard begin marker indicates a standard private key string; leave it unchanged
	if !strings.HasPrefix(c.PrivateKey, privateBeginStr) {
		//Read from the file if it exists
		if fileExists(c.PrivateKey) {
			contents, err := fileGetContents(c.PrivateKey)
			if err == nil {
				c.PrivateKey = string(contents)
			}
		}
	}

	//Remove the prefix and suffix
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

// Load loads configuration
//
// Load and parse the public and private keys if either is present
func (c *Conf) Load() {
	var (
		publicBeginStr  = constants.PublicKeyBeginStr
		publicEndStr    = constants.PublicKeyEndStr
		privateBeginStr = constants.PrivateKeyBeginStr
		privateEndStr   = constants.PrivateKeyEndStr
	)

	if len(c.PublicKey) != 0 {
		//A standard begin marker indicates a standard public key string; leave it unchanged
		if !strings.HasPrefix(c.PublicKey, publicBeginStr) {
			//Read from the file if it exists
			if fileExists(c.PublicKey) {
				contents, err := fileGetContents(c.PublicKey)
				if err == nil {
					c.PublicKey = string(contents)
				}
			}
		}

		//Remove the prefix and suffix
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
		//A standard begin marker indicates a standard private key string; leave it unchanged
		if !strings.HasPrefix(c.PrivateKey, privateBeginStr) {
			//Read from the file if it exists
			if fileExists(c.PrivateKey) {
				contents, err := fileGetContents(c.PrivateKey)
				if err == nil {
					c.PrivateKey = string(contents)
				}
			}
		}

		//Remove the prefix and suffix
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

// GetRSAPrivateKeyObj retrieves the RSA private key object
func (c *Conf) GetRSAPrivateKeyObj() *rsa.PrivateKey {
	return c.rsaPrivateKey
}

// GetRSAPublicKeyObj retrieves the RSA public key object
func (c *Conf) GetRSAPublicKeyObj() *rsa.PublicKey {
	return c.rsaPublicKey
}

// GetED25519PrivateKeyObj retrieves the Ed25519 private key object
func (c *Conf) GetED25519PrivateKeyObj() crypto.PrivateKey {
	return c.ed25519PrivateKey
}

// GetED25519PublicKeyObj retrieves the Ed25519 public key object
func (c *Conf) GetED25519PublicKeyObj() crypto.PublicKey {
	return c.ed25519PublicKey
}

// GetECDSAPrivateKeyObj retrieves the ECDSA private key object
func (c *Conf) GetECDSAPrivateKeyObj() *ecdsa.PrivateKey {
	return c.ecdsaPrivateKey
}

// GetECDSAPublicKeyObj retrieves the ECDSA public key object
func (c *Conf) GetECDSAPublicKeyObj() *ecdsa.PublicKey {
	return c.ecdsaPublicKey
}
