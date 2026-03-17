package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {
	t.Run("Bcrypt", func(t *testing.T) {
		assert.NotNil(t, Bcrypt())
	})
}

func TestBcryptHashMake(t *testing.T) {
	t.Run("BcryptHashMake", func(t *testing.T) {
		result, err := Bcrypt().HashMake("password")
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
	t.Run("BcryptHashMake-Cost-overflow", func(t *testing.T) {
		result, err := Bcrypt().HashMake("password", 1)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
	t.Run("BcryptHashMake-Cost", func(t *testing.T) {
		result, err := Bcrypt().HashMake("password", 4)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestBcryptHashCheck(t *testing.T) {
	t.Run("BcryptHashCheck", func(t *testing.T) {
		hashed, err := Bcrypt().HashMake("password")
		assert.Nil(t, err)
		assert.NotEmpty(t, hashed)
		ok, err := Bcrypt().HashCheck("password", hashed)
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}
