package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"unicode/utf8"
)

type stringImpl struct {
}

type IString interface {
	// Wordwrap 以给定的字符和长度来打断字符串
	Wordwrap(rawStr string, length int, split string) string
	// WrapRedisKey 包装redis键名
	//
	// 给redis的键加入应用名前缀，如：
	//
	// appName=game key=user
	//
	// 最终的redis键名为：game:user
	//
	// 此方法的主要作用是按应用来划分redis键名
	WrapRedisKey(appName, key string) string
	// RandomLetters 随机字符串(字母)
	RandomLetters(length int) string
	// RandomDigitalChars 随机字符串(数字)
	RandomDigitalChars(length int) string
	// RandomString 随机字符串(字母+数字)
	RandomString(length int) string
	// RandomComplexString 随机字符串(可带特殊符号)
	RandomComplexString(length int) string
	// Reverse 翻转字符串
	Reverse(s string) string
	// Shuffle 打乱字符串
	Shuffle(s string) string
	// PaddingLeft 向左填充字符串
	//
	// rawString 原字符
	//
	// padChar 填充字符
	//
	// length 最终字符长度
	PaddingLeft(rawString, padChar string, length int) string
	// PaddingRight 向右填充字符串
	//
	// rawString 原字符
	//
	// padChar 填充字符
	//
	// length 最终字符长度
	PaddingRight(rawString, padChar string, length int) string
	// PaddingBoth 向两端填充字符串
	//
	// rawString 原字符
	//
	// padChar 填充字符
	//
	// length 最终字符长度
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
	PaddingBoth(rawString, padChar string, length int) string
	// FromCharCode 返回ASCII码对应的字符
	//
	// # Note
	//
	// 常规ASCII码表范围为0~127
	//
	// 扩展ASCII码表范围为128~255
	//
	// more: https://www.rfc-editor.org/rfc/rfc698.txt
	FromCharCode(code int32) string
	// CharCodeAt 返回字符对应的ASCII码
	CharCodeAt(character string) rune
	//Truncate 按指定长度截取字符串
	//
	// length若超出字符串长度则返回字符串本身
	Truncate(rawString string, length int) string
}

var _ IString = &stringImpl{}

// String 实例化string工具类
func String() IString {
	return &stringImpl{}
}

// Wordwrap 以给定的字符和长度来打断字符串
func (stringImpl) Wordwrap(rawStr string, length int, split string) string {
	if len(rawStr) <= length || length < 1 {
		return rawStr
	}

	strSplit := strings.Split(rawStr, "")

	var (
		start    int
		end      = length
		finalStr = strings.Builder{}
	)

	for {
		if start > len(strSplit) {
			break
		}
		if end >= len(strSplit) {
			finalStr.WriteString(strings.Join(strSplit[start:], ""))
		} else {
			finalStr.WriteString(strings.Join(strSplit[start:end], ""))
			finalStr.WriteString(split)
		}

		start = end
		end += length
	}

	return finalStr.String()
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
func (stringImpl) WrapRedisKey(appName, key string) string {
	return fmt.Sprintf("%s:%s", appName, key)
}

const (
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitalChars    = "0123456789"
	specificSymbols = "~`!@#$%^&*-_+=/\\|/<>.:;'\""
)

// RandomLetters 随机字符串(字母)
func (stringImpl) RandomLetters(length int) string {
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
func (stringImpl) RandomDigitalChars(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = digitalChars[rand.Intn(len(digitalChars))]
	}

	return string(b)
}

// RandomString 随机字符串(字母+数字)
func (stringImpl) RandomString(length int) string {
	var s = fmt.Sprintf("%s%s", letters, digitalChars)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}

	return string(b)
}

// RandomComplexString 随机字符串(可带特殊符号)
func (stringImpl) RandomComplexString(length int) string {
	var s = fmt.Sprintf("%s%s%s", letters, digitalChars, specificSymbols)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}

	return string(b)
}

// Reverse 翻转字符串
func (stringImpl) Reverse(s string) string {
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

// Shuffle 打乱字符串
func (stringImpl) Shuffle(s string) string {
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

// PaddingLeft 向左填充字符串
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
func (stringImpl) PaddingLeft(rawString, padChar string, length int) string {
	return padding(rawString, padChar, length, 0)
}

// PaddingRight 向右填充字符串
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
func (stringImpl) PaddingRight(rawString, padChar string, length int) string {
	return padding(rawString, padChar, length, 1)
}

// PaddingBoth 向两端填充字符串
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
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
func (stringImpl) PaddingBoth(rawString, padChar string, length int) string {
	return padding(rawString, padChar, length, 2)
}

// padding 填充字符串
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
//
// padType 0:向左填充,1:向右填充,2:向两端填充
func padding(rawString, padChar string, length, padType int) string {
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

// FromCharCode 返回ASCII码对应的字符
//
// # Note
//
// 常规ASCII码表范围为0~127
//
// 扩展ASCII码表范围为128~255
//
// more: https://www.rfc-editor.org/rfc/rfc698.txt
func (stringImpl) FromCharCode(code int32) string {
	return fmt.Sprintf("%c", code)
}

// CharCodeAt 返回字符对应的ASCII码
func (stringImpl) CharCodeAt(character string) rune {
	return ([]rune(character))[0]
}

// Truncate 按指定长度截取字符串
//
// length若超出字符串长度则返回字符串本身
func (stringImpl) Truncate(rawString string, length int) string {
	if length <= 0 {
		return ""
	}
	if len(rawString) == 0 {
		return rawString
	}
	if utf8.RuneCountInString(rawString) <= length {
		return rawString
	}
	return string([]rune(rawString)[:length])
}
