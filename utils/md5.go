package utils

import (
	"crypto/md5"
	"encoding/hex"
)

type md5impl struct {
}

type IMD5 interface {
	// Encode md5编码
	Encode(rawString string) string
}

var _ IMD5 = &md5impl{}

// MD5 实例化md5工具类
func MD5() IMD5 {
	return &md5impl{}
}

// Encode md5编码
func (md5impl) Encode(rawString string) string {
	instance := md5.New()
	_, _ = instance.Write([]byte(rawString))
	sumString := instance.Sum(nil)
	return hex.EncodeToString(sumString)
}
