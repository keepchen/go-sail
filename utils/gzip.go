package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

// GzipCompress gzip压缩
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
func GzipDecompress(content []byte) []byte {
	r, e := gzip.NewReader(bytes.NewReader(content))
	if e != nil {
		return []byte(``)
	}
	res, _ := io.ReadAll(r)
	_ = r.Close()

	return res
}
