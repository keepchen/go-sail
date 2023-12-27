package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5Encrypt md5加密
func MD5Encrypt(rawString string) string {
	instance := md5.New()
	_, _ = instance.Write([]byte(rawString))
	sumString := instance.Sum(nil)
	return hex.EncodeToString(sumString)
}
