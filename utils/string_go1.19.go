//go:build go1.19 && !go1.22
// +build go1.19,!go1.22

package utils

import (
	"reflect"
	"unsafe"
)

// StrToBytes converts a string to a byte slice
//
// nolint
func StrToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

// StrToBytes converts a string to a byte slice
//
// nolint
func (stringImpl) StrToBytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{Data: sh.Data, Len: sh.Len, Cap: sh.Len}

	return *(*[]byte)(unsafe.Pointer(&bh))
}

// BytesToStr converts a byte slice to a string
//
// nolint
func BytesToStr(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}

	return *(*string)(unsafe.Pointer(&sh))
}

// BytesToStr converts a byte slice to a string
//
// nolint
func (stringImpl) BytesToStr(b []byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{Data: bh.Data, Len: bh.Len}

	return *(*string)(unsafe.Pointer(&sh))
}
