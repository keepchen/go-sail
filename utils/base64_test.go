package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Encode(t *testing.T) {
	encodedString := Base64Encode([]byte(rawString))
	t.Log(encodedString)
	assert.Equal(t, "aGVsbG8gd29ybGQh", encodedString)
}

func TestBase64Decode(t *testing.T) {
	encodedString := Base64Encode([]byte(rawString))

	decodeBytes, err := Base64Decode(encodedString)
	t.Log(rawString, encodedString, string(decodeBytes))
	assert.NoError(t, err)
	assert.Equal(t, string(decodeBytes), rawString)
}
