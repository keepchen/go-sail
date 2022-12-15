package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var rawString = "hello world!"

func TestAesEncode(t *testing.T) {
	encodedString, err := AesEncode(rawString)
	t.Log(encodedString)
	assert.NoError(t, err)
}

func TestAesDecode(t *testing.T) {
	encodedString, err := AesEncode(rawString)
	t.Log(encodedString)
	assert.NoError(t, err)
	decodedString, err := AesDecode(encodedString)
	t.Log(decodedString)
	assert.NoError(t, err)
	assert.Equal(t, decodedString, rawString)
}
