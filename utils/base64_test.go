package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64ImplEncode(t *testing.T) {
	encodedString := Base64().Encode([]byte(rawString))
	t.Log(encodedString)
	assert.Equal(t, "aGVsbG8gd29ybGQh", encodedString)
}

func TestBase64ImplDecode(t *testing.T) {
	encodedString := Base64().Encode([]byte(rawString))

	decodeBytes, err := Base64().Decode(encodedString)
	t.Log(rawString, encodedString, string(decodeBytes))
	assert.NoError(t, err)
	assert.Equal(t, string(decodeBytes), rawString)
}
