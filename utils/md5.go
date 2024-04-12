package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encode md5编码
func MD5Encode(rawString string) string {
	instance := md5.New()
	_, _ = instance.Write([]byte(rawString))
	sumString := instance.Sum(nil)
	return hex.EncodeToString(sumString)
}
