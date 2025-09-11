package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

type gzipImpl struct {
}

type IGzip interface {
	// Compress gzip压缩
	//
	// 使用 gzip.DefaultCompression 压缩等级
	Compress(content []byte) []byte
	// Decompress gzip解压
	Decompress(content []byte) []byte
}

var gi IGzip = &gzipImpl{}

// Gzip 实例化gzip工具类
func Gzip() IGzip {
	return gi
}

// Compress gzip压缩
//
// 使用 gzip.DefaultCompression 压缩等级
func (gzipImpl) Compress(content []byte) []byte {
	var b bytes.Buffer
	w, e := gzip.NewWriterLevel(&b, gzip.DefaultCompression)
	if e != nil {
		return content
	}

	_, _ = w.Write(content)
	_ = w.Close()

	return b.Bytes()
}

// Decompress gzip解压
func (gzipImpl) Decompress(content []byte) []byte {
	r, e := gzip.NewReader(bytes.NewReader(content))
	if e != nil {
		return []byte(``)
	}
	res, _ := io.ReadAll(r)
	_ = r.Close()

	return res
}
