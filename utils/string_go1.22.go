//go:build go1.22
// +build go1.22

package utils

import "unsafe"

// StrToBytes 字符串转换为字节数组
func StrToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// BytesToStr 字节数组转换为字符串
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
