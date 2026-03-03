package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipImplCompressAndDecompress(t *testing.T) {
	var bytesArr = [][]byte{
		[]byte("hello world!"),
		[]byte("hello golang!"),
	}

	t.Run("TestGzipImplCompress", func(t *testing.T) {
		for _, byt := range bytesArr {
			assert.NotEmpty(t, Gzip().Compress(byt))
		}
	})

	t.Run("TestGzipImplDecompress", func(t *testing.T) {
		for _, byt := range bytesArr {
			s := Gzip().Compress(byt)
			assert.NotEmpty(t, string(Gzip().Decompress(s)))
		}
	})
}
