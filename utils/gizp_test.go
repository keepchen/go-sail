package utils

import "testing"

func TestGzipCompressAndDecompress(t *testing.T) {
	var bytesArr = [][]byte{
		[]byte("hello world!"),
		[]byte("hello golang!"),
	}

	t.Run("TestGzipCompress", func(t *testing.T) {
		for _, byt := range bytesArr {
			t.Log(GzipCompress(byt))
		}
	})

	t.Run("TestGzipDecompress", func(t *testing.T) {
		for _, byt := range bytesArr {
			s := GzipCompress(byt)
			t.Log(string(s), "--->", string(GzipDecompress(s)))
		}
	})
}
