package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawString = "hello world!"

func TestAesEncode(t *testing.T) {
	encodedString, err := AesEncode(rawString, KEY)
	t.Log(encodedString)
	assert.NoError(t, err)
}

func TestAesDecode(t *testing.T) {
	encodedString, err := AesEncode(rawString, KEY)
	t.Log(encodedString)
	assert.NoError(t, err)
	decodedString, err := AesDecode(encodedString, KEY)
	t.Log(decodedString)
	assert.NoError(t, err)
	assert.Equal(t, decodedString, rawString)
}
