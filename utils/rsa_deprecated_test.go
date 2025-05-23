package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	privateKeyDeprecated = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDUvUDx+LPQ0S+L
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
RIXvUK2duHjDxiaPKtANi2P4
-----END PRIVATE KEY-----
`)
	publicKeyDeprecated = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1L1A8fiz0NEvi/uVJrQ9
hCcNS/ed5lQjFgSZLQW6j0yIQ5q3d9/t3DatLe617gC4aKmUv+d5w3IVNcLAr3on
jQ0e0rHTuBzEuG58YMrmApAtboAHcrydhQr+8VKaXOz78iENYEEv5zj9MCnGYhYq
saVQaFSs0r/08OHLfKf1i9ru+cWpGY+aDOs7U3cGMMZ0XLORMpVmdUA+5v8vacCz
u/Wq722YIEH39RpL2ydxTP2q5bBNrLeDWdxyBxZlP9zZ/T7AW+boqIGDdPHKBmWr
KTJQ+GGzUqOGzruYQ5sM3TnU8Avb4OF36uyADBwA4bP944tKSNSET7BC3N0UerRo
QwIDAQAB
-----END PUBLIC KEY-----`)
)

func TestRSAEncode(t *testing.T) {
	encodedString, err := RSAEncrypt(rawString, publicKeyDeprecated)
	t.Log(encodedString)
	assert.NoError(t, err)
}

func TestRSADecode(t *testing.T) {
	t.Run("Decode", func(t *testing.T) {
		encodedString, err := RSAEncrypt(rawString, publicKeyDeprecated)
		t.Log(encodedString)
		assert.NoError(t, err)
		decodedString, err := RSADecrypt(encodedString, privateKeyDeprecated)
		t.Log(decodedString)
		assert.NoError(t, err)
		assert.Equal(t, decodedString, rawString)
	})

	t.Run("Decode-Error", func(t *testing.T) {
		encodedString, err := RSAEncrypt(rawString, publicKeyDeprecated)
		t.Log(encodedString)
		assert.NoError(t, err)
		decodedString, err := RSADecrypt("", privateKeyDeprecated)
		t.Log(decodedString)
		assert.Error(t, err)
		assert.NotEqual(t, decodedString, rawString)
	})

	t.Run("Decode-Error2", func(t *testing.T) {
		assert.Panics(t, func() {
			encodedString, err := RSAEncrypt(rawString, publicKeyDeprecated)
			t.Log(encodedString)
			assert.NoError(t, err)
			decodedString, err := RSADecrypt(encodedString, []byte(strings.Replace(string(privateKeyDeprecated), "a", "b", 1)))
			t.Log(decodedString)
			assert.Error(t, err)
			assert.NotEqual(t, decodedString, rawString)
		})
	})
}

func TestRSASign(t *testing.T) {
	t.Run("Sign", func(t *testing.T) {
		sign, err := RSASign([]byte(rawString), privateKeyDeprecated)
		t.Log(Base64Encode([]byte(sign)))
		assert.NoError(t, err)
	})

	t.Run("Sign-Error", func(t *testing.T) {
		assert.Panics(t, func() {
			sign, err := RSASign([]byte(rawString), []byte(strings.Replace(string(privateKeyDeprecated), "a", "b", 1)))
			t.Log(Base64Encode([]byte(sign)))
			assert.Error(t, err)
		})
	})
}

func TestRSAVerifySign(t *testing.T) {
	t.Run("VerifySign", func(t *testing.T) {
		sign, err := RSASign([]byte(rawString), privateKeyDeprecated)
		t.Log(sign)
		assert.NoError(t, err)

		ok, err := RSAVerifySign([]byte(rawString), []byte(sign), publicKeyDeprecated)
		assert.NoError(t, err)
		assert.Equal(t, true, ok)
	})

	t.Run("VerifySign-Error", func(t *testing.T) {
		sign, err := RSASign([]byte(rawString), privateKeyDeprecated)
		t.Log(sign)
		assert.NoError(t, err)

		ok, err := RSAVerifySign([]byte(rawString), []byte(sign+"a"), publicKeyDeprecated)
		assert.Error(t, err)
		assert.NotEqual(t, true, ok)
	})

	t.Run("VerifySign-Error2", func(t *testing.T) {
		sign, err := RSASign([]byte(rawString), privateKeyDeprecated)
		t.Log(sign)
		assert.NoError(t, err)

		ok, err := RSAVerifySign([]byte(rawString), []byte(sign), []byte(strings.Replace(string(publicKeyDeprecated), "a", "b", 1)))
		assert.Error(t, err)
		assert.NotEqual(t, true, ok)
	})
}
