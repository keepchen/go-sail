package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

// GzipCompress gzip压缩
//
// Deprecated: GzipCompress is deprecated,it will be removed in the future.
//
// Please use Gzip().Compress() instead.
//
// 使用 gzip.DefaultCompression 压缩等级
func GzipCompress(content []byte) []byte {
	var b bytes.Buffer
	w, e := gzip.NewWriterLevel(&b, gzip.DefaultCompression)
	if e != nil {
		return content
	}

	_, _ = w.Write(content)
	_ = w.Close()

	return b.Bytes()
}

// GzipDecompress gzip解压
//
// Deprecated: GzipDecompress is deprecated,it will be removed in the future.
//
// Please use Gzip().Decompress() instead.
func GzipDecompress(content []byte) []byte {
	r, e := gzip.NewReader(bytes.NewReader(content))
	if e != nil {
		return []byte(``)
	}
	res, _ := io.ReadAll(r)
	_ = r.Close()

	return res
}
