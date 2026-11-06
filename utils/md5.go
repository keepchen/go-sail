package utils

import (
	"crypto/md5"
	"encoding/hex"
)

type md5impl struct {
}

type IMD5 interface {
	// Encode MD5 encoding
	Encode(rawString string) string
}

var mdi IMD5 = &md5impl{}

// MD5 instantiates MD5 utilities
func MD5() IMD5 {
	return mdi
}

// Encode MD5 encoding
func (md5impl) Encode(rawString string) string {
	instance := md5.New()
	_, _ = instance.Write([]byte(rawString))
	sumString := instance.Sum(nil)
	return hex.EncodeToString(sumString)
}
