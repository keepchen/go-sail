package etcd

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("New-NonValue", func(t *testing.T) {
		conf := Conf{}
		t.Log(New(conf))
	})

	t.Run("New", func(t *testing.T) {
		conf := Conf{
			Tls:      &tls.Config{},
			Username: "username",
			Password: "password",
		}
		t.Log(New(conf))
	})
}

func TestInit(t *testing.T) {
	t.Run("Init-NonValue", func(t *testing.T) {
		conf := Conf{}
		assert.Panics(t, func() {
			Init(conf)
		})
	})

	t.Run("Init", func(t *testing.T) {
		conf := Conf{
			Tls:      &tls.Config{},
			Username: "username",
			Password: "password",
		}
		assert.Panics(t, func() {
			Init(conf)
		})
	})
}
