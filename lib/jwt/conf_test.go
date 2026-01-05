package jwt

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/constants"

	"github.com/stretchr/testify/assert"
)

func TestMustLoad(t *testing.T) {
	t.Run("MustLoad-NonValue", func(t *testing.T) {
		conf := Conf{
			Enable: true,
		}
		conf.MustLoad()
	})

	t.Run("MustLoad-FromFile", func(t *testing.T) {
		var (
			publicKeyFile  = fmt.Sprintf("./tester-publicKey-%d", time.Now().UnixMilli())
			privateKeyFile = fmt.Sprintf("./tester-privateKey-%d", time.Now().UnixMilli())
		)

		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  publicKeyFile,
			PrivateKey: privateKeyFile,
		}

		err := os.WriteFile(publicKeyFile, rsaPublicKey, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(privateKeyFile, rsaPrivateKey, 0644)
		assert.NoError(t, err)

		conf.MustLoad()

		err = os.Remove(publicKeyFile)
		assert.NoError(t, err)
		err = os.Remove(privateKeyFile)
		assert.NoError(t, err)
	})

	t.Run("MustLoad-NonPrefixOrSuffix", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  strings.Replace(string(rsaPublicKey), constants.PublicKeyBeginStr, "", 1),
			PrivateKey: strings.Replace(string(rsaPrivateKey), constants.PrivateKeyBeginStr, "", 1),
		}

		conf.MustLoad()
	})

	t.Run("MustLoad", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  string(rsaPublicKey),
			PrivateKey: string(rsaPrivateKey),
		}
		conf.MustLoad()
	})
}

func TestLoad(t *testing.T) {
	t.Run("Load-NonValue", func(t *testing.T) {
		conf := Conf{
			Enable: true,
		}
		conf.Load()
	})

	t.Run("Load-FromFile", func(t *testing.T) {
		var (
			publicKeyFile  = fmt.Sprintf("./tester-publicKey-%d", time.Now().UnixMilli())
			privateKeyFile = fmt.Sprintf("./tester-privateKey-%d", time.Now().UnixMilli())
		)

		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  publicKeyFile,
			PrivateKey: privateKeyFile,
		}

		err := os.WriteFile(publicKeyFile, rsaPublicKey, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(privateKeyFile, rsaPrivateKey, 0644)
		assert.NoError(t, err)

		conf.Load()

		err = os.Remove(publicKeyFile)
		assert.NoError(t, err)
		err = os.Remove(privateKeyFile)
		assert.NoError(t, err)
	})

	t.Run("Load-NonPrefixOrSuffix", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  strings.Replace(string(rsaPublicKey), constants.PublicKeyBeginStr, "", 1),
			PrivateKey: strings.Replace(string(rsaPrivateKey), constants.PrivateKeyBeginStr, "", 1),
		}

		conf.Load()
	})

	t.Run("Load", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			PublicKey:  string(rsaPublicKey),
			PrivateKey: string(rsaPrivateKey),
		}
		conf.Load()
	})
}

func TestGetRSAPrivateKeyObj(t *testing.T) {
	t.Run("GetRSAPrivateKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  string(rsaPublicKey),
			PrivateKey: string(rsaPrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetRSAPrivateKeyObj())
	})
}

func TestGetRSAPublicKeyObj(t *testing.T) {
	t.Run("GetRSAPublicKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodRS256.String(),
			PublicKey:  string(rsaPublicKey),
			PrivateKey: string(rsaPrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetRSAPublicKeyObj())
	})
}

func TestGetED25519PrivateKeyObj(t *testing.T) {
	t.Run("GetED25519PrivateKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodEdDSA.String(),
			PublicKey:  string(ed25519PublicKey),
			PrivateKey: string(ed25519PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetED25519PrivateKeyObj())
	})
}

func TestGetED25519PublicKeyObj(t *testing.T) {
	t.Run("GetED25519PublicKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodEdDSA.String(),
			PublicKey:  string(ed25519PublicKey),
			PrivateKey: string(ed25519PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetED25519PublicKeyObj())
	})
}

func TestGetECDSAPrivateKeyObj(t *testing.T) {
	t.Run("GetECDSA256PrivateKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodES256.String(),
			PublicKey:  string(ecdsa256PublicKey),
			PrivateKey: string(ecdsa256PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetECDSAPrivateKeyObj())
	})

	t.Run("GetECDSA384PrivateKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodES384.String(),
			PublicKey:  string(ecdsa384PublicKey),
			PrivateKey: string(ecdsa384PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetECDSAPrivateKeyObj())
	})

	t.Run("GetECDSA521PrivateKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodES512.String(),
			PublicKey:  string(ecdsa521PublicKey),
			PrivateKey: string(ecdsa521PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetECDSAPrivateKeyObj())
	})
}

func TestGetECDSAPublicKeyObj(t *testing.T) {
	t.Run("GetECDSA256PublicKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodES256.String(),
			PublicKey:  string(ecdsa256PublicKey),
			PrivateKey: string(ecdsa256PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetECDSAPublicKeyObj())
	})

	t.Run("GetECDSA384PublicKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodES384.String(),
			PublicKey:  string(ecdsa384PublicKey),
			PrivateKey: string(ecdsa384PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetECDSAPublicKeyObj())
	})

	t.Run("GetECDSA521PublicKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			Algorithm:  SigningMethodES512.String(),
			PublicKey:  string(ecdsa521PublicKey),
			PrivateKey: string(ecdsa521PrivateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetECDSAPublicKeyObj())
	})
}
