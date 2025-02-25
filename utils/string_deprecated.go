package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

// Wordwrap 以给定的字符和长度来打断字符串
//
// Deprecated: Wordwrap is deprecated,it will be removed in the future.
//
// Please use String().Wordwrap() instead.
func Wordwrap(rawStr string, length int, split string) string {
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
// Deprecated: WrapRedisKey is deprecated,it will be removed in the future.
//
// Please use String().WrapRedisKey() instead.
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
	lettersDeprecated         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitalCharsDeprecated    = "0123456789"
	specificSymbolsDeprecated = "~`!@#$%^&*-_+=/\\|/<>.:;'\""
)

// RandomLetters 随机字符串(字母)
//
// Deprecated: RandomLetters is deprecated,it will be removed in the future.
//
// Please use String().RandomLetters() instead.
func RandomLetters(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = lettersDeprecated[rand.Intn(len(lettersDeprecated))]
	}

	return string(b)
}

// RandomDigitalChars 随机字符串(数字)
//
// Deprecated: RandomDigitalChars is deprecated,it will be removed in the future.
//
// Please use String().RandomDigitalChars() instead.
func RandomDigitalChars(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = digitalCharsDeprecated[rand.Intn(len(digitalCharsDeprecated))]
	}

	return string(b)
}

// RandomString 随机字符串(字母+数字)
//
// Deprecated: RandomString is deprecated,it will be removed in the future.
//
// Please use String().RandomString() instead.
func RandomString(length int) string {
	var s = fmt.Sprintf("%s%s", lettersDeprecated, digitalCharsDeprecated)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}

	return string(b)
}

// RandomComplexString 随机字符串(可带特殊符号)
//
// Deprecated: RandomComplexString is deprecated,it will be removed in the future.
//
// Please use String().RandomComplexString() instead.
func RandomComplexString(length int) string {
	var s = fmt.Sprintf("%s%s%s", lettersDeprecated, digitalCharsDeprecated, specificSymbolsDeprecated)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}

	return string(b)
}

// StringReverse 翻转字符串
//
// Deprecated: StringReverse is deprecated,it will be removed in the future.
//
// Please use String().Reverse() instead.
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
//
// Deprecated: StringShuffle is deprecated,it will be removed in the future.
//
// Please use String().Shuffle() instead.
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
// Deprecated: StringPaddingLeft is deprecated,it will be removed in the future.
//
// Please use String().PaddingLeft() instead.
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
func StringPaddingLeft(rawString, padChar string, length int) string {
	return paddingStringDeprecated(rawString, padChar, length, 0)
}

// StringPaddingRight 向右填充字符串
//
// Deprecated: StringPaddingRight is deprecated,it will be removed in the future.
//
// Please use String().PaddingRight() instead.
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
func StringPaddingRight(rawString, padChar string, length int) string {
	return paddingStringDeprecated(rawString, padChar, length, 1)
}

// StringPaddingBoth 向两端填充字符串
//
// Deprecated: StringPaddingBoth is deprecated,it will be removed in the future.
//
// Please use String().PaddingBoth() instead.
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
func StringPaddingBoth(rawString, padChar string, length int) string {
	return paddingStringDeprecated(rawString, padChar, length, 2)
}

// paddingStringDeprecated 填充字符串
//
// rawString 原字符
//
// padChar 填充字符
//
// length 最终字符长度
//
// padType 0:向左填充,1:向右填充,2:向两端填充
func paddingStringDeprecated(rawString, padChar string, length, padType int) string {
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
// Deprecated: FromCharCode is deprecated,it will be removed in the future.
//
// Please use String().FromCharCode() instead.
//
// # Note
//
// 常规ASCII码表范围为0~127
//
// 扩展ASCII码表范围为128~255
//
// more: https://www.rfc-editor.org/rfc/rfc698.txt
func FromCharCode(code int32) string {
	return fmt.Sprintf("%c", code)
}

// CharCodeAt 返回字符对应的ASCII码
//
// Deprecated: CharCodeAt is deprecated,it will be removed in the future.
//
// Please use String().CharCodeAt() instead.
func CharCodeAt(character string) rune {
	return ([]rune(character))[0]
}
