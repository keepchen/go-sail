package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawString = "hello world!"

// KeyForTest 测试密钥
const KeyForTest = "fakeKeyChangeMe!"

func TestAesImplEncode(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		encodedString, err := Aes().Encode(rawString, KeyForTest)
		assert.NoError(t, err)
		assert.NotEmpty(t, encodedString)
	})

	t.Run("Encode-Error", func(t *testing.T) {
		encodedString, err := Aes().Encode(rawString, "")
		assert.Error(t, err)
		assert.Equal(t, "", encodedString)
	})
}

func TestAesImplDecode(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		encodedString, err := Aes().Encode(rawString, KeyForTest)
		assert.NoError(t, err)
		decodedString, err := Aes().Decode(encodedString, KeyForTest)
		assert.NoError(t, err)
		assert.Equal(t, decodedString, rawString)
	})

	t.Run("Decode-Error", func(t *testing.T) {
		decodedString, err := Aes().Decode(rawString, KeyForTest)
		assert.Error(t, err)
		assert.NotEqual(t, decodedString, rawString)
	})
}
