package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawStringDeprecated = "hello world!"

// KeyForTestDeprecated 测试密钥
const KeyForTestDeprecated = "fakeKeyChangeMe!"

func TestAesEncode(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		encodedString, err := AesEncode(rawStringDeprecated, KeyForTestDeprecated)
		t.Log(encodedString)
		assert.NoError(t, err)
	})

	t.Run("Encode-Error", func(t *testing.T) {
		encodedString, err := AesEncode(rawStringDeprecated, "")
		t.Log(encodedString)
		assert.Error(t, err)
	})
}

func TestAesDecode(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		encodedString, err := AesEncode(rawStringDeprecated, KeyForTestDeprecated)
		t.Log(encodedString)
		assert.NoError(t, err)
		decodedString, err := AesDecode(encodedString, KeyForTestDeprecated)
		t.Log(decodedString)
		assert.NoError(t, err)
		assert.Equal(t, decodedString, rawStringDeprecated)
	})

	t.Run("Decode-Error", func(t *testing.T) {
		decodedString, err := AesDecode(rawString, KeyForTestDeprecated)
		t.Log(decodedString)
		assert.Error(t, err)
		assert.NotEqual(t, decodedString, rawStringDeprecated)
	})
}
