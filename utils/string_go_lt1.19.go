//go:build !go1.19
// +build !go1.19

package utils

func (stringImpl) StrToBytes(s string) []byte {
	return []byte(s)
}

func (stringImpl) BytesToStr(b []byte) string {
	return string(b)
}
