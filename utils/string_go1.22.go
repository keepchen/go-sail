//go:build go1.22
// +build go1.22

package utils

import "unsafe"

// StrToBytes converts a string to a byte slice
func StrToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// StrToBytes converts a string to a byte slice
func (stringImpl) StrToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}

// BytesToStr converts a byte slice to a string
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// BytesToStr converts a byte slice to a string
func (stringImpl) BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
