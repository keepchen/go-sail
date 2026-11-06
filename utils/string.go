package utils

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand/v2"
	"strings"
	"unicode/utf8"
)

type stringImpl struct {
}

type IString interface {
	// Wordwrap wraps a string using the given character and length
	Wordwrap(rawStr string, length int, split string) string
	// WrapRedisKey wraps a Redis key name
	//
	// Prefix Redis keys with the application name, for example:
	//
	// appName=game key=user
	//
	// The resulting Redis key is: game:user
	//
	// This method primarily separates Redis key names by application
	WrapRedisKey(appName, key string) string
	// RandomLetters random string (letters)
	RandomLetters(length int) string
	// RandomDigitalChars random string (digits)
	RandomDigitalChars(length int) string
	// RandomString random string (letters+digits)
	RandomString(length int) string
	// RandomComplexString random string (may include special characters)
	RandomComplexString(length int) string
	// Reverse reverses a string
	Reverse(s string) string
	// Shuffle shuffles a string
	Shuffle(s string) string
	// PaddingLeft pads a string on the left
	//
	// rawString original string
	//
	// padChar padding character
	//
	// length final string length
	PaddingLeft(rawString, padChar string, length int) string
	// PaddingRight pads a string on the right
	//
	// rawString original string
	//
	// padChar padding character
	//
	// length final string length
	PaddingRight(rawString, padChar string, length int) string
	// PaddingBoth pads both ends of a string
	//
	// rawString original string
	//
	// padChar padding character
	//
	// length final string length
	//
	// # Note
	//
	// If the padding cannot be split evenly, add one extra character on the right, for example:
	//
	// rawString = "a",padChar = "#",length = 4
	//
	// Then:
	//
	// result = "#a##"
	PaddingBoth(rawString, padChar string, length int) string
	// FromCharCode returns the character for an ASCII code
	//
	// # Note
	//
	// The standard ASCII range is 0~127
	//
	// The extended ASCII range is 128~255
	//
	// more: https://www.rfc-editor.org/rfc/rfc698.txt
	FromCharCode(code int32) string
	// CharCodeAt returns the ASCII code for a character
	CharCodeAt(character string) rune
	//Truncate truncates a string to the specified length
	//
	// Return the original string if length exceeds its length
	Truncate(rawString string, length int) string
	// StrToBytes converts a string to a byte slice
	StrToBytes(s string) []byte
	// BytesToStr converts a byte slice to a string
	BytesToStr(b []byte) string
	// RandomIdentity generates random IDs
	//
	// Generate a random string of the specified length, excluding ambiguous characters such as 1, l, and I, for unpredictable identifiers with collision constraints,
	// such as filenames and random account names.
	//
	// Since len(alphabet) = 32
	//
	// Total space = 32^length = 2^(5*length)
	//
	// [Collision probability]
	//
	// Approximate formula: P ≈ n² / (2N)
	//
	//   - n=total number generated
	//
	//   - N=total space
	//
	// Practical example, length=10
	//
	//	N=2^50≈1.12e15
	//
	//	n=1e5
	//
	// P ≈ (1e5)^2 / (2 × 1.12e15) ≈ 1e10 / 2.24e15 ≈ 4.5e-6 ≈ 0.00045%
	//
	// length = 12, P ≈ 0.00000043%
	RandomIdentity(length int) (string, error)
	// MaskString masks sensitive information in a string
	//
	// rawString - original string
	//
	// maskChar - masking character
	//
	// maskLen - number of masking characters
	//
	// prefixLen - length of the original prefix to retain
	//
	// suffixLen - length of the original suffix to retain
	//
	// Example 1:
	//
	// rawString = 123456789
	//
	// markChar = *
	//
	// maskLen = 4
	//
	// prefixLen = 2
	//
	// suffixLen = 2
	//
	// Result: 12****89
	//
	// ------------------------------------------
	//
	// Example 2:
	//
	// rawString = 12
	//
	// markChar = *
	//
	// maskLen = 4
	//
	// prefixLen = 2
	//
	// suffixLen = 2
	//
	// Result: 12****
	MaskString(rawString string, maskChar string, maskLen, prefixLen, suffixLen int) string
	// RandHexString generates a hexadecimal string of the specified length
	//
	// length - string length
	RandHexString(length int) string
}

var sti IString = &stringImpl{}

// String instantiates string utilities
func String() IString {
	return sti
}

// Wordwrap wraps a string using the given character and length
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

// WrapRedisKey wraps a Redis key name
//
// Prefix Redis keys with the application name, for example:
//
// appName=game key=user
//
// The resulting Redis key is: game:user
//
// This method primarily separates Redis key names by application
func (stringImpl) WrapRedisKey(appName, key string) string {
	return fmt.Sprintf("%s:%s", appName, key)
}

const (
	letters         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitalChars    = "0123456789"
	specificSymbols = "~`!@#$%^&*-_+=/\\|/<>.:;'\""
)

// RandomLetters random string (letters)
func (stringImpl) RandomLetters(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}

	return string(b)
}

// RandomDigitalChars random string (digits)
func (stringImpl) RandomDigitalChars(length int) string {
	if length < 1 {
		return ""
	}

	b := make([]byte, length)

	for i := range b {
		b[i] = digitalChars[rand.IntN(len(digitalChars))]
	}

	return string(b)
}

// RandomString random string (letters+digits)
func (stringImpl) RandomString(length int) string {
	var s = fmt.Sprintf("%s%s", letters, digitalChars)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.IntN(len(s))]
	}

	return string(b)
}

// RandomComplexString random string (may include special characters)
func (stringImpl) RandomComplexString(length int) string {
	var s = fmt.Sprintf("%s%s%s", letters, digitalChars, specificSymbols)

	b := make([]byte, length)

	for i := range b {
		b[i] = s[rand.IntN(len(s))]
	}

	return string(b)
}

// Reverse reverses a string
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

// Shuffle shuffles a string
func (stringImpl) Shuffle(s string) string {
	length := len(s)

	if length < 1 {
		return s
	}

	arr := strings.Split(s, "")

	for range arr {
		r0, r1 := rand.IntN(length), rand.IntN(length)
		arr[r0], arr[r1] = arr[r1], arr[r0]
	}

	return strings.Join(arr, "")
}

// PaddingLeft pads a string on the left
//
// rawString original string
//
// padChar padding character
//
// length final string length
func (stringImpl) PaddingLeft(rawString, padChar string, length int) string {
	return padding(rawString, padChar, length, 0)
}

// PaddingRight pads a string on the right
//
// rawString original string
//
// padChar padding character
//
// length final string length
func (stringImpl) PaddingRight(rawString, padChar string, length int) string {
	return padding(rawString, padChar, length, 1)
}

// PaddingBoth pads both ends of a string
//
// rawString original string
//
// padChar padding character
//
// length final string length
//
// # Note
//
// If the padding cannot be split evenly, add one extra character on the right, for example:
//
// rawString = "a",padChar = "#",length = 4
//
// Then:
//
// result = "#a##"
func (stringImpl) PaddingBoth(rawString, padChar string, length int) string {
	return padding(rawString, padChar, length, 2)
}

// padding pads a string
//
// rawString original string
//
// padChar padding character
//
// length final string length
//
// padType 0:pad left,1:pad right,2:pad both ends
func padding(rawString, padChar string, length, padType int) string {
	if length < 1 || len(padChar) == 0 {
		return rawString
	}
	padLength := length - len(rawString)
	if padLength < 1 {
		return rawString
	}

	//Use the first character if the padding character length exceeds 1
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
		//If the padding cannot be split evenly, add one extra character on the right
		if padLength&1 == 1 {
			right += 1
		}
		return fmt.Sprintf("%s%s%s", strings.Repeat(padChar, left), rawString, strings.Repeat(padChar, right))
	}
}

// FromCharCode returns the character for an ASCII code
//
// # Note
//
// The standard ASCII range is 0~127
//
// The extended ASCII range is 128~255
//
// more: https://www.rfc-editor.org/rfc/rfc698.txt
func (stringImpl) FromCharCode(code int32) string {
	return fmt.Sprintf("%c", code)
}

// CharCodeAt returns the ASCII code for a character
func (stringImpl) CharCodeAt(character string) rune {
	return ([]rune(character))[0]
}

// Truncate truncates a string to the specified length
//
// Return the original string if length exceeds its length
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

const alphabet = "23456789abcdefghjkmnpqrstuvwxyz"

// RandomIdentity generates random IDs
//
// Generate a random string of the specified length, excluding ambiguous characters such as 1, l, and I, for unpredictable identifiers with collision constraints,
// such as filenames and random account names.
//
// Since len(alphabet) = 32
//
// Total space = 32^length = 2^(5*length)
//
// [Collision probability]
//
// Approximate formula: P ≈ n² / (2N)
//
//   - n=total number generated
//
//   - N=total space
//
// Practical example, length=10
//
//	N=2^50≈1.12e15
//
//	n=1e5
//
// P ≈ (1e5)^2 / (2 × 1.12e15) ≈ 1e10 / 2.24e15 ≈ 4.5e-6 ≈ 0.00045%
//
// length = 12, P ≈ 0.00000043%
func (stringImpl) RandomIdentity(length int) (string, error) {
	b := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := crand.Int(crand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}

	return string(b), nil
}

// MaskString masks sensitive information in a string
//
// rawString - original string
//
// maskChar - masking character
//
// maskLen - number of masking characters
//
// prefixLen - length of the original prefix to retain
//
// suffixLen - length of the original suffix to retain
//
// Example 1:
//
// rawString = 123456789
//
// markChar = *
//
// maskLen = 4
//
// prefixLen = 2
//
// suffixLen = 2
//
// Result: 12****89
//
// ------------------------------------------
//
// Example 2:
//
// rawString = 12
//
// markChar = *
//
// maskLen = 4
//
// prefixLen = 2
//
// suffixLen = 2
//
// Result: 12****
func (stringImpl) MaskString(rawString string, maskChar string, maskLen, prefixLen, suffixLen int) string {
	if len(rawString) == 0 {
		return ""
	}

	runes := []rune(rawString)
	n := len(runes)

	// Boundary checks
	if prefixLen < 0 {
		prefixLen = 0
	}
	if suffixLen < 0 {
		suffixLen = 0
	}
	if maskLen < 1 {
		maskLen = 4
	}

	var builder strings.Builder

	// Prefix (at most n characters)
	if prefixLen > n {
		prefixLen = n
	}
	builder.WriteString(string(runes[:prefixLen]))

	// Force maskLen padding characters
	builder.WriteString(strings.Repeat(maskChar, maskLen))

	// Suffix (must not overlap the prefix)
	if suffixLen > 0 {
		start := n - suffixLen
		if start < prefixLen {
			start = prefixLen
		}
		builder.WriteString(string(runes[start:]))
	}

	return builder.String()
}

// RandHexString generates a hexadecimal string of the specified length
//
// length - string length
func (stringImpl) RandHexString(length int) string {
	if length <= 0 {
		return ""
	}

	// hex uses 1 byte -> 2 chars
	byteLen := (length + 1) / 2

	b := make([]byte, byteLen)
	if _, err := crand.Read(b); err != nil {
		return ""
	}

	s := hex.EncodeToString(b)

	// Truncate to the specified length
	return s[:length]
}
