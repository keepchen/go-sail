package jwt

import (
	"fmt"
	"github.com/keepchen/go-sail/v3/constants"
	"os"
	"strings"
	"testing"
	"time"

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
			PublicKey:  publicKeyFile,
			PrivateKey: privateKeyFile,
		}

		err := os.WriteFile(publicKeyFile, publicKey, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(privateKeyFile, privateKey, 0644)
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
			PublicKey:  strings.Replace(string(publicKey), constants.PublicKeyBeginStr, "", 1),
			PrivateKey: strings.Replace(string(privateKey), constants.PrivateKeyBeginStr, "", 1),
		}

		conf.MustLoad()
	})

	t.Run("MustLoad", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			PublicKey:  string(publicKey),
			PrivateKey: string(privateKey),
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
			PublicKey:  publicKeyFile,
			PrivateKey: privateKeyFile,
		}

		err := os.WriteFile(publicKeyFile, publicKey, 0644)
		assert.NoError(t, err)
		err = os.WriteFile(privateKeyFile, privateKey, 0644)
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
			PublicKey:  strings.Replace(string(publicKey), constants.PublicKeyBeginStr, "", 1),
			PrivateKey: strings.Replace(string(privateKey), constants.PrivateKeyBeginStr, "", 1),
		}

		conf.Load()
	})

	t.Run("Load", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			PublicKey:  string(publicKey),
			PrivateKey: string(privateKey),
		}
		conf.Load()
	})
}

func TestGetPrivateKeyObj(t *testing.T) {
	t.Run("GetPrivateKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			PublicKey:  string(publicKey),
			PrivateKey: string(privateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetPrivateKeyObj())
	})
}

func TestGetPublicKeyObj(t *testing.T) {
	t.Run("GetPublicKeyObj", func(t *testing.T) {
		conf := Conf{
			Enable:     true,
			PublicKey:  string(publicKey),
			PrivateKey: string(privateKey),
		}
		conf.MustLoad()
		assert.NotNil(t, conf.GetPublicKeyObj())
	})
}
