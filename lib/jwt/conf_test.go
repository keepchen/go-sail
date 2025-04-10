package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustLoad(t *testing.T) {
	conf := Conf{
		Enable:     true,
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}
	t.Run("MustLoad", func(t *testing.T) {
		conf.MustLoad()
	})
}

func TestLoad(t *testing.T) {
	conf := Conf{
		Enable:     true,
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}
	t.Run("Load", func(t *testing.T) {
		conf.Load()
	})
}

func TestGetPrivateKeyObj(t *testing.T) {
	conf := Conf{
		Enable:     true,
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}
	conf.MustLoad()
	t.Run("GetPrivateKeyObj", func(t *testing.T) {
		assert.NotNil(t, conf.GetPrivateKeyObj())
	})
}

func TestGetPublicKeyObj(t *testing.T) {
	conf := Conf{
		Enable:     true,
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}
	conf.MustLoad()
	t.Run("GetPublicKeyObj", func(t *testing.T) {
		assert.NotNil(t, conf.GetPublicKeyObj())
	})
}
