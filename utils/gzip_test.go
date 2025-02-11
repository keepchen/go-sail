package utils

import "testing"

func TestGzipImplCompressAndDecompress(t *testing.T) {
	var bytesArr = [][]byte{
		[]byte("hello world!"),
		[]byte("hello golang!"),
	}

	t.Run("TestGzipImplCompress", func(t *testing.T) {
		for _, byt := range bytesArr {
			t.Log(Gzip().Compress(byt))
		}
	})

	t.Run("TestGzipImplDecompress", func(t *testing.T) {
		for _, byt := range bytesArr {
			s := Gzip().Compress(byt)
			t.Log(string(s), "--->", string(Gzip().Decompress(s)))
		}
	})
}
