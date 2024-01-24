package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

// Wordwrap 以给定的字符和长度来打断字符串
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

// WrapRedisKey 包装redis键名
//
// 给redis的键加入应用名前缀，如：
//
// appName=game key=user
//
// 最终的redis键名为：game:user
//
// 此方法的主要作用是按应用来划分redis键名
func WrapRedisKey(appName, key string) string {
	return fmt.Sprintf("%s:%s", appName, key)
}

const (
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitalChars    = "0123456789"
	specificSymbols = "~`!@#$%^&*-_+=/\\|/<>.:;'\""
)

// RandomLetters 随机字符串(字母)
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

// RandomDigitalChars 随机字符串(数字)
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

// RandomComplexString 随机字符串(可带特殊符号)
func RandomComplexString(length int) string {
	var s = fmt.Sprintf("%s%s%s", letters, digitalChars, specificSymbols)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}

	return string(b)
}

// StringReverse 翻转字符串
func StringReverse(s string) string {
	length := len(s)

	if length < 1 {
		return s
	}

	b := make([]byte, length)

	for i := range b {
		b[length-1-i] = s[i]
	}

	return string(b)
}

// StringShuffle 打乱字符串
func StringShuffle(s string) string {
	length := len(s)

	if length < 1 {
		return s
	}

	arr := strings.Split(s, "")

	for range arr {
		r0, r1 := rand.Intn(length), rand.Intn(length)
		arr[r0], arr[r1] = arr[r1], arr[r0]
	}

	return strings.Join(arr, "")
}

// StringPaddingLeft 向左填充字符串
//
// @param rawString 原字符
//
// @param padChar 填充字符
//
// @param length 最终字符长度
func StringPaddingLeft(rawString, padChar string, length int) string {
	return paddingString(rawString, padChar, length, 0)
}

// StringPaddingRight 向右填充字符串
//
// @param rawString 原字符
//
// @param padChar 填充字符
//
// @param length 最终字符长度
func StringPaddingRight(rawString, padChar string, length int) string {
	return paddingString(rawString, padChar, length, 1)
}

// StringPaddingBoth 向两端填充字符串
//
// @param rawString 原字符
//
// @param padChar 填充字符
//
// @param length 最终字符长度
//
// # Note
//
// 如果填充长度不能均分，那么右侧多填充一个字符，如：
//
// rawString = "a",padChar = "#",length = 4
//
// 则：
//
// result = "#a##"
func StringPaddingBoth(rawString, padChar string, length int) string {
	return paddingString(rawString, padChar, length, 2)
}

// paddingString 填充字符串
//
// @param rawString 原字符
//
// @param padChar 填充字符
//
// @param length 最终字符长度
//
// @param padType 0:向左填充,1:向右填充,2:向两端填充
func paddingString(rawString, padChar string, length, padType int) string {
	if length < 1 || len(padChar) == 0 {
		return rawString
	}
	padLength := length - len(rawString)
	if padLength < 1 {
		return rawString
	}

	//如果填充字符长度大于1，则取第一个字符
	if len(padChar) > 1 {
		var s = make([]byte, 1)
		s[0] = padChar[0]
		padChar = string(s)
	}

	switch padType {
	default:
		return fmt.Sprintf("%s%s", strings.Repeat(padChar, padLength), rawString)
	case 0:
		return fmt.Sprintf("%s%s", strings.Repeat(padChar, padLength), rawString)
	case 1:
		return fmt.Sprintf("%s%s", rawString, strings.Repeat(padChar, padLength))
	case 2:
		left, right := padLength/2, padLength/2
		//如果填充长度不能均分，那么右侧多填充一个字符
		if padLength&1 == 1 {
			right += 1
		}
		return fmt.Sprintf("%s%s%s", strings.Repeat(padChar, left), rawString, strings.Repeat(padChar, right))
	}
}
