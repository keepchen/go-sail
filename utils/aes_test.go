package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawString = "hello world!"

// KeyForTest 测试密钥
const KeyForTest = "fakeKeyChangeMe!"

func TestAesImplEncode(t *testing.T) {
	encodedString, err := Aes().Encode(rawString, KeyForTest)
	t.Log(encodedString)
	assert.NoError(t, err)
}

func TestAesImplDecode(t *testing.T) {
	encodedString, err := Aes().Encode(rawString, KeyForTest)
	t.Log(encodedString)
	assert.NoError(t, err)
	decodedString, err := Aes().Decode(encodedString, KeyForTest)
	t.Log(decodedString)
	assert.NoError(t, err)
	assert.Equal(t, decodedString, rawString)
}
