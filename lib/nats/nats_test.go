package nats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	t.Run("GetInstance", func(t *testing.T) {
		assert.Nil(t, GetInstance())
	})
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		assert.Panics(t, func() {
			conf := Conf{}
			Init(conf)
		})
	})

	t.Run("Init-Credentials", func(t *testing.T) {
		assert.Panics(t, func() {
			conf := Conf{
				Username: "username",
				Password: "password",
			}
			Init(conf)
		})
	})
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		conf := Conf{}
		t.Log(New(conf))
	})

	t.Run("New-Credentials", func(t *testing.T) {
		conf := Conf{
			Username: "username",
			Password: "password",
		}
		t.Log(New(conf))
	})
}
