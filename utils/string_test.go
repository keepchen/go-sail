package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var rawStrArr = []string{
	"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",
	"abc",
	"",
}

var expectedStrArr = []string{
	`MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDUvUDx+LPQ0S+L
+5UmtD2EJw1L953mVCMWBJktBbqPTIhDmrd33+3cNq0t7rXuALhoqZS/53nDchU1
wsCveieNDR7SsdO4HMS4bnxgyuYCkC1ugAdyvJ2FCv7xUppc7PvyIQ1gQS/nOP0w
KcZiFiqxpVBoVKzSv/Tw4ct8p/WL2u75xakZj5oM6ztTdwYwxnRcs5EylWZ1QD7m
/y9pwLO79arvbZggQff1GkvbJ3FM/arlsE2st4NZ3HIHFmU/3Nn9PsBb5uiogYN0
8coGZaspMlD4YbNSo4bOu5hDmwzdOdTwC9vg4Xfq7IAMHADhs/3ji0pI1IRPsELc
3RR6tGhDAgMBAAECggEAQal4VjcxKQ6n4kjwrFWNdzCmhgATmHf3rGAW9zKBdqFk
nZkvb6yKOiIWKcs4FBHc2VEePG0xxAV+Tm2iE4dclciq7tU8R+N5RIO1mBqIC9p8
a1LQ+bUF2X6fWdTpGC19Riq1ejQkmPWaEDeUp8m3u8UOoGUiQppE++R1bjBZNaT5
S16qbfDOV9plF550wnwbq6fNZlWT4PdiI6ox/4KPZdhIKGKnKkh4xX4mHk3E9fl+
udHbXiT3qjSDOEchUpHglzNZG1LMD2BWb+zcxUbJzm2r5BviZmKHPd4+w5mt+kfP
bDHnFCgjnlZoFswNFO2s/ZHk99NveoDa1i0OVGbuwQKBgQD045x067oKAvBcRawt
P5H6DaifBFnwp3xTg8GIPwitD8+bQMIi6s0jbN7HI7A7S5BTFnukCweURhBiXmbr
s29ImiaxcIVGdBKDhkHABq3oino4oHVs8bLpD9moaQccon79aPQ38j8KeU1UJ7a4
R6Jbd2eLoZPmj5bPrQHkbqDWWwKBgQDeZDjS0tLtOjbY9CHdB9+BWSyL2DylqSye
c/Ew/c8sr2SBK59Db0W2Vgc62iTOYjzWYBTBUrYWRRoAnSoQLkePjQ+mpzGMtpR9
BKq3ADrIremgJGRFIo+NL9qjpJu238na+FGp1DgfTSXMMxzLC23lTgh6PXIIcF03
kL/yN9lqOQKBgQDPCmCMuX9gV3u/h2g6GTThpAqb5qHjxLZoJUzKVACR0HxFVkrM
Gpe1C6aN1q54czpiBPAjkO+nfFT91bJONDYxu6JbAjarihbc+/U61GrT37/VgFPG
99G7GZt7ttA8dWXH+aQAaN7DjCrEq47f3jB2BE2Wz9SraVqn2i1vY9i3YQKBgQCu
0d4RbIU+0upWteMA27WI+s6XyA40s75NeRr6xipcGCxLlj0GR6xnX00jqGQSkQr+
Al2OczSMYRnFrcZpHdhHMj5BZWEAGm6zsD16ygVrx7rFlpXz+u0ZsaqPxVBa+6S0
K0wW0qqjgIPb97oEqyFihmsHnNHNbHb6vSEGiXyxkQKBgAwb/3lWqp1Zpj6hMw9N
dB0c6huQYLqX2INkKj9PcIlFq0nOeHMZfMisuQKhvcGsPQsHMP2NbPjZiLnbpRHP
vplU0p7ayaXuNF2t73k/L5f92+8VBuYECEUOXw2xST5gvkPdKGK1xM1cLT6y8TrF
RIXvUK2duHjDxiaPKtANi2P4`,
	`MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1L1A8fiz0NEvi/uVJrQ9
hCcNS/ed5lQjFgSZLQW6j0yIQ5q3d9/t3DatLe617gC4aKmUv+d5w3IVNcLAr3on
jQ0e0rHTuBzEuG58YMrmApAtboAHcrydhQr+8VKaXOz78iENYEEv5zj9MCnGYhYq
saVQaFSs0r/08OHLfKf1i9ru+cWpGY+aDOs7U3cGMMZ0XLORMpVmdUA+5v8vacCz
u/Wq722YIEH39RpL2ydxTP2q5bBNrLeDWdxyBxZlP9zZ/T7AW+boqIGDdPHKBmWr
KTJQ+GGzUqOGzruYQ5sM3TnU8Avb4OF36uyADBwA4bP944tKSNSET7BC3N0UerRo
QwIDAQAB`,
	"abc",
	"",
}

func TestStringImplWordwrap(t *testing.T) {
	for i := 0; i < len(rawStrArr); i++ {
		finalStr := String().Wordwrap(rawStrArr[i], 64, "\n")
		//t.Log(finalStr, "<---->", expectedStrArr[i])
		assert.Equal(t, expectedStrArr[i], finalStr)
	}
}

func TestStringImplWrapRedisKey(t *testing.T) {
	holders := [][]string{
		{"game", "user", "game:user"},
		{"video", "drama", "video:drama"},
		{"video", "drama:romantic", "video:drama:romantic"},
	}

	for _, arr := range holders {
		key := String().WrapRedisKey(arr[0], arr[1])
		//t.Log(key)
		assert.Equal(t, arr[2], key)
	}
}

func TestStringImplRandomLetters(t *testing.T) {
	holders := []int{0, 1, 3, 5, 7, 8, 9, 33}
	for _, v := range holders {
		s := String().RandomLetters(v)
		//t.Log(s)
		assert.Equal(t, v, len(s))
	}
}

func TestStringImplRandomDigitalChars(t *testing.T) {
	holders := []int{0, 1, 3, 5, 7, 8, 9, 33}
	for _, v := range holders {
		s := String().RandomDigitalChars(v)
		//t.Log(s)
		assert.Equal(t, v, len(s))
	}
}

func TestStringImplRandomString(t *testing.T) {
	holders := []int{0, 1, 3, 5, 7, 8, 9, 33, 100}
	for _, v := range holders {
		s := String().RandomString(v)
		//t.Log(s)
		assert.Equal(t, v, len(s))
	}
}

func TestStringImplRandomComplexString(t *testing.T) {
	holders := []int{0, 1, 3, 5, 7, 8, 9, 33, 100}
	for _, v := range holders {
		s := String().RandomComplexString(v)
		//t.Log(s)
		assert.Equal(t, v, len(s))
	}
}

func TestStringImplStringReverse(t *testing.T) {
	holders := []string{"", "a", "ab", "abc", "1234567890"}
	result := []string{"", "a", "ba", "cba", "0987654321"}
	for index, v := range holders {
		s := String().Reverse(v)
		//t.Log(s)
		assert.Equal(t, result[index], s)
	}
}

func TestStringImplStringShuffle(t *testing.T) {
	holders := []string{"", "a", "ab", "abc", "1234567890"}
	for _, v := range holders {
		s := String().Shuffle(v)
		t.Log(s)
	}
}

func TestStringImplPaddingLeft(t *testing.T) {
	holders := []struct {
		raw     string
		padChar string
		length  int
		result  string
	}{
		{"", "=", 1, "="},
		{"", "=", 2, "=="},
		{"", "=", 3, "==="},
		{"a", "=", 1, "a"},
		{"a", "=", 2, "=a"},
		{"a", "=", 3, "==a"},
	}
	for _, v := range holders {
		s := String().PaddingLeft(v.raw, v.padChar, v.length)
		//t.Log(s)
		assert.Equal(t, v.result, s)
	}
}

func TestStringImplPaddingRight(t *testing.T) {
	holders := []struct {
		raw     string
		padChar string
		length  int
		result  string
	}{
		{"", "=", 1, "="},
		{"", "=", 2, "=="},
		{"", "=", 3, "==="},
		{"a", "=", 1, "a"},
		{"a", "=", 2, "a="},
		{"a", "=", 3, "a=="},
	}
	for _, v := range holders {
		s := String().PaddingRight(v.raw, v.padChar, v.length)
		//t.Log(s)
		assert.Equal(t, v.result, s)
	}
}

func TestStringImplPaddingBoth(t *testing.T) {
	holders := []struct {
		raw     string
		padChar string
		length  int
		result  string
	}{
		{"", "=", 1, "="},
		{"", "=", 2, "=="},
		{"", "=", 3, "==="},
		{"a", "=", 1, "a"},
		{"a", "=", 2, "a="},
		{"a", "+=", 3, "+a+"},
		{"a", "=+", 4, "=a=="},
	}
	for _, v := range holders {
		s := String().PaddingBoth(v.raw, v.padChar, v.length)
		//t.Log(s)
		assert.Equal(t, v.result, s)
	}
}

var a2zLowercase = "abcdefghijklmnopqrstuvwxyz"

func TestStringImplCharCodeAt(t *testing.T) {
	arr := strings.Split(a2zLowercase, "")
	asciiLowercase := int32(97)
	asciiUppercase := int32(65)
	for _, letter := range arr {
		//lowercase
		code := String().CharCodeAt(letter)
		//t.Log(letter, ":", code)
		assert.Equal(t, code, asciiLowercase)
		asciiLowercase++

		//uppercase
		upLetter := strings.ToUpper(letter)
		codeUp := String().CharCodeAt(upLetter)
		//t.Log(upLetter, ":", codeUp)
		assert.Equal(t, codeUp, asciiUppercase)
		asciiUppercase++
	}
}

func TestStringImplFromCharCode(t *testing.T) {
	arr := strings.Split(a2zLowercase, "")
	asciiLowercase := int32(97)
	asciiUppercase := int32(65)
	for _, letter := range arr {
		//lowercase
		character := String().FromCharCode(asciiLowercase)
		//t.Log(asciiLowercase, ":", character)
		assert.Equal(t, letter, character)
		asciiLowercase++

		//uppercase
		upLetter := strings.ToUpper(letter)
		characterUp := String().FromCharCode(asciiUppercase)
		//t.Log(asciiUppercase, ":", characterUp)
		assert.Equal(t, upLetter, characterUp)
		asciiUppercase++
	}
}

func TestStringImplCharCodeRange(t *testing.T) {
	for i := 0; i < 256; i++ {
		character := String().FromCharCode(int32(i))
		code := String().CharCodeAt(character)
		//t.Log("character:", character, "code:", code, []byte(character), int32(i))
		assert.Equal(t, int32(i), code)
	}
}

func TestStringImplBytesStrExchange(t *testing.T) {
	t.Run("TestBytesStrExchange", func(t *testing.T) {
		for _, str := range expectedStrArr {
			assert.Equal(t, str, BytesToStr(StrToBytes(str)))
		}
	})
}

func TestStringImplTruncate(t *testing.T) {
	holders := []string{"", "a", "ab", "abc", "1234567890", "abc123你好こんにちは안녕하세요"}
	for _, v := range holders {
		for i := -1; i < 20; i++ {
			s := String().Truncate(v, i)
			t.Log(s)
		}
	}
}
