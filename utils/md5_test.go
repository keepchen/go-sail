package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5ImplEncrypt(t *testing.T) {
	encodedString := MD5().Encode(rawString)
	t.Log(encodedString)
	assert.Equal(t, "fc3ff98e8c6a0d3087d515c0473f8677", encodedString)
}
