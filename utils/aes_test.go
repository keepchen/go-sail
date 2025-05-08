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
		t.Log(encodedString)
		assert.NoError(t, err)
	})

	t.Run("Encode-Error", func(t *testing.T) {
		encodedString, err := Aes().Encode(rawString, "")
		t.Log(encodedString)
		assert.Error(t, err)
	})
}

func TestAesImplDecode(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		encodedString, err := Aes().Encode(rawString, KeyForTest)
		t.Log(encodedString)
		assert.NoError(t, err)
		decodedString, err := Aes().Decode(encodedString, KeyForTest)
		t.Log(decodedString)
		assert.NoError(t, err)
		assert.Equal(t, decodedString, rawString)
	})

	t.Run("Decode-Error", func(t *testing.T) {
		decodedString, err := Aes().Decode(rawString, KeyForTest)
		t.Log(decodedString)
		assert.Error(t, err)
		assert.NotEqual(t, decodedString, rawString)
	})
}
