package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var anotherArgon2Config = Argon2Config{
	Memory:      32 * 1024, //32MB
	Iterations:  6,
	Parallelism: 8,
	SaltLength:  16,
	KeyLength:   32,
}

func TestArgon2(t *testing.T) {
	t.Run("Argon2", func(t *testing.T) {
		assert.NotNil(t, Argon2())
	})
	t.Run("Argon2-With-Config", func(t *testing.T) {
		assert.NotNil(t, Argon2(anotherArgon2Config))
	})
}

func TestArgon2HashMake(t *testing.T) {
	t.Run("Argon2HashMake", func(t *testing.T) {
		result, err := Argon2().HashMake("password")
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
	t.Run("Argon2HashMake-Another-Config", func(t *testing.T) {
		result, err := Argon2(anotherArgon2Config).HashMake("password")
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestArgon2HashCheck(t *testing.T) {
	t.Run("Argon2HashCheck", func(t *testing.T) {
		hashed, err := Argon2().HashMake("password")
		assert.Nil(t, err)
		assert.NotEmpty(t, hashed)
		valid, err := Argon2().HashCheck("password", hashed)
		assert.Nil(t, err)
		assert.True(t, valid)
	})
	t.Run("Argon2HashCheck-Another-Config", func(t *testing.T) {
		hashed, err := Argon2(anotherArgon2Config).HashMake("password")
		assert.Nil(t, err)
		assert.NotEmpty(t, hashed)
		valid, err := Argon2().HashCheck("password", hashed)
		assert.Nil(t, err)
		assert.True(t, valid)
	})
}
