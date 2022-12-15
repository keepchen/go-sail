package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

//Wordwrap 以给定的字符和长度来打断字符串
func Wordwrap(rawStr string, length int, split string) string {
	if len(rawStr) <= length || length < 1 {
		return rawStr
	}

	strSplit := strings.Split(rawStr, "")

	var (
		start    int
		end      = length
		finalStr string
	)

	for {
		if start > len(strSplit) {
			break
		}
		if end >= len(strSplit) {
			finalStr += strings.Join(strSplit[start:], "")
		} else {
			finalStr += strings.Join(strSplit[start:end], "") + split
		}

		start = end
		end += length
	}

	return finalStr
}

//WrapRedisKey 包装redis键名
//
//给redis的键加入应用名前缀，如：
//
//appName=game key=user
//
//最终的redis键名为：game:user
//
//此方法的主要作用是按应用来划分redis键名
func WrapRedisKey(appName, key string) string {
	return fmt.Sprintf("%s:%s", appName, key)
}

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitalChars = "0123456789"
)

//RandomLetters 随机字符串(字母)
func RandomLetters(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

//RandomDigitalChars 随机字符串(数字)
func RandomDigitalChars(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = digitalChars[rand.Intn(len(digitalChars))]
	}

	return string(b)
}
