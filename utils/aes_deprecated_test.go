package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawStringDeprecated = "hello world!"

// KeyForTestDeprecated 测试密钥
const KeyForTestDeprecated = "fakeKeyChangeMe!"

func TestAesEncode(t *testing.T) {
	encodedString, err := AesEncode(rawStringDeprecated, KeyForTestDeprecated)
	t.Log(encodedString)
	assert.NoError(t, err)
}

func TestAesDecode(t *testing.T) {
	encodedString, err := AesEncode(rawStringDeprecated, KeyForTestDeprecated)
	t.Log(encodedString)
	assert.NoError(t, err)
	decodedString, err := AesDecode(encodedString, KeyForTestDeprecated)
	t.Log(decodedString)
	assert.NoError(t, err)
	assert.Equal(t, decodedString, rawStringDeprecated)
}
