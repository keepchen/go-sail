//go:build go1.19 && !go1.22
// +build go1.19,!go1.22

package utils

import (
	"reflect"
	"unsafe"
)

// StrToBytes 字符串转换为字节数组
//
// nolint
func StrToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

// BytesToStr 字节数组转换为字符串
//
// nolint
func BytesToStr(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}

	return *(*string)(unsafe.Pointer(&sh))
}
