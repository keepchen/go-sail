package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

var (
	rsaPrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
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

	rsaPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1L1A8fiz0NEvi/uVJrQ9
hCcNS/ed5lQjFgSZLQW6j0yIQ5q3d9/t3DatLe617gC4aKmUv+d5w3IVNcLAr3on
jQ0e0rHTuBzEuG58YMrmApAtboAHcrydhQr+8VKaXOz78iENYEEv5zj9MCnGYhYq
saVQaFSs0r/08OHLfKf1i9ru+cWpGY+aDOs7U3cGMMZ0XLORMpVmdUA+5v8vacCz
u/Wq722YIEH39RpL2ydxTP2q5bBNrLeDWdxyBxZlP9zZ/T7AW+boqIGDdPHKBmWr
KTJQ+GGzUqOGzruYQ5sM3TnU8Avb4OF36uyADBwA4bP944tKSNSET7BC3N0UerRo
QwIDAQAB
-----END PUBLIC KEY-----`)

	ed25519PublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAULc3WGEBwX9ZnS5n44dWwf9/4rZOuDWNvSoxHw4nPUc=
-----END PUBLIC KEY-----`)

	ed25519PrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIGsoHvm67+edKGCXSa+LQO7mUWzA1cFF/FB7B2cc2DWj
-----END PRIVATE KEY-----`)

	ecdsa256PublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmCPXCEvaUaLR682mRKIWYtpS2l7O
V1UN75b5zuEz8MdK8UEtcJpvzsYkbYaK001ekjjCYvWTvFWZJqazdm2KnQ==
-----END PUBLIC KEY-----`)

	ecdsa256PrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgJXd7bPHBJfP0b6MS
KWEaAJJMu6TzMfpeL/a+1Nx116yhRANCAASYI9cIS9pRotHrzaZEohZi2lLaXs5X
VQ3vlvnO4TPwx0rxQS1wmm/OxiRthorTTV6SOMJi9ZO8VZkmprN2bYqd
-----END PRIVATE KEY-----`)

	ecdsa384PublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEp8drJIb1rXHmeI0HQoA5ObnC92FhskT6
4n49e4YArW7i1XXPql08CZQ56zRbCJGxRAIkJMatYcWNJUPtHWzsKXSUtjYydONx
ZMrEV9gp+Ql8/7cjB2EWVNOtN9KjTVn9
-----END PUBLIC KEY-----`)

	ecdsa384PrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDD+m7bfdzLjT+sOGbAX
uIhYI54ll1gZSaCr44E81XDfLkw4JqTmkezxRKnMwPIIgdmhZANiAASnx2skhvWt
ceZ4jQdCgDk5ucL3YWGyRPrifj17hgCtbuLVdc+qXTwJlDnrNFsIkbFEAiQkxq1h
xY0lQ+0dbOwpdJS2NjJ043FkysRX2Cn5CXz/tyMHYRZU06030qNNWf0=
-----END PRIVATE KEY-----`)

	ecdsa521PublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBLCS7Dyu2SpLBiTwHk8oo9aCsy1f4
WN9h1lhtXw/8F0g5QLiOi5fmesE6y8cYg4vFn3R11191Okmb9Cxm39mSmGUAh6MJ
V3Wx390urz5Dmzk+VjKA5P6Du/mPlm3B/X9cVT0rVIfYVDY5gwc68a5m7ehggkjJ
Sg+Ux9Qh1+hxxbQfALA=
-----END PUBLIC KEY-----`)

	ecdsa521PrivateKey = []byte(`-----BEGIN PRIVATE KEY-----
MIHuAgEAMBAGByqGSM49AgEGBSuBBAAjBIHWMIHTAgEBBEIAzIkcHJlsIQlvIMqE
b9TQua6LaxRfr2NI2vSgWXW/KSM0YJbIu9aoI4eTKOEjqIOowq5MbHw0U8VmiY0d
9q2VYVqhgYkDgYYABAEsJLsPK7ZKksGJPAeTyij1oKzLV/hY32HWWG1fD/wXSDlA
uI6Ll+Z6wTrLxxiDi8WfdHXXX3U6SZv0LGbf2ZKYZQCHowlXdbHf3S6vPkObOT5W
MoDk/oO7+Y+WbcH9f1xVPStUh9hUNjmDBzrxrmbt6GCCSMlKD5TH1CHX6HHFtB8A
sA==
-----END PRIVATE KEY-----`)
)

var (
	confArr = []Conf{
		{
			Algorithm:  SigningMethodHS512.String(),
			HmacSecret: "a-b-c-d-e-f-g",
		},
		{
			Algorithm:  SigningMethodRS256.String(),
			PrivateKey: string(rsaPrivateKey),
			PublicKey:  string(rsaPublicKey),
		},
		{
			Algorithm:  SigningMethodRS512.String(),
			PrivateKey: string(rsaPrivateKey),
			PublicKey:  string(rsaPublicKey),
		},
		{
			Algorithm:  SigningMethodEdDSA.String(),
			PrivateKey: string(ed25519PrivateKey),
			PublicKey:  string(ed25519PublicKey),
		},
		{
			Algorithm:  SigningMethodES256.String(),
			PrivateKey: string(ecdsa256PrivateKey),
			PublicKey:  string(ecdsa256PublicKey),
		},
		{
			Algorithm:  SigningMethodES384.String(),
			PrivateKey: string(ecdsa384PrivateKey),
			PublicKey:  string(ecdsa384PublicKey),
		},
		{
			Algorithm:  SigningMethodES512.String(),
			PrivateKey: string(ecdsa521PrivateKey),
			PublicKey:  string(ecdsa521PublicKey),
		},
	}

	appClaim = AppClaims{
		Name: "test",
		RegisteredClaims: jwtLib.RegisteredClaims{
			ExpiresAt: jwtLib.NewNumericDate(time.Unix(time.Now().Unix()+1500, 0)),
			Issuer:    defaultTokenIssuer,
		},
	}

	mapClaims = []MapClaims{
		{
			"name":  "test",
			"aud":   "test-user",
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"iat":   time.Now().Unix(),
			"iss":   "go-sail",
			"nbf":   time.Now().Add(10 * time.Minute).Unix(),
			"sub":   "tester",
			"valid": false,
		},
		{
			"name":  "test",
			"aud":   "test-user",
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
			"iat":   time.Now().Unix(),
			"iss":   "go-sail",
			"nbf":   time.Now().Add(-10 * time.Minute).Unix(),
			"sub":   "tester",
			"valid": true,
		},
	}
)

func TestSign(t *testing.T) {
	t.Run("Sign", func(t *testing.T) {
		for _, conf := range confArr {
			t.Log("----------------", conf.Algorithm, "----------------")
			conf.Load()
			token, err := Sign(appClaim, conf)
			t.Log(err)
			t.Log("struct claim:", token)
			assert.NoError(t, err)
			for _, claim := range mapClaims {
				token2, err := SignWithMap(claim, conf)
				t.Log(err)
				t.Log("map claim:", token2)
				assert.NoError(t, err)
			}
		}
	})

	t.Run("Sign-NonValue", func(t *testing.T) {
		for _, conf := range confArr {
			t.Log("----------------", conf.Algorithm, "----------------")
			conf.Load()
			appClaim.Issuer = ""
			token, err := Sign(appClaim, conf)
			t.Log(err)
			t.Log("struct claim:", token)
			assert.NoError(t, err)
			for _, claim := range mapClaims {
				token2, err := SignWithMap(claim, conf)
				t.Log(err)
				t.Log("map claim:", token2)
				assert.NoError(t, err)
			}
		}
	})
}

func TestVerify(t *testing.T) {
	t.Run("Verify", func(t *testing.T) {
		for _, conf := range confArr {
			t.Log("----------------", conf.Algorithm, "----------------")
			conf.Load()
			token, err := Sign(appClaim, conf)
			t.Log(conf.Algorithm, err)
			t.Log("signed claim:", token)
			assert.NoError(t, err)

			claim, err := Verify(token, conf)
			t.Log(conf.Algorithm, err)
			t.Log("verify claim:", claim)
			assert.NoError(t, err)

			for _, claim := range mapClaims {
				delete(claim, "iss")
				token2, err := SignWithMap(claim, conf)
				t.Log(conf.Algorithm, err)
				t.Log("map signed token:", token2)
				assert.NoError(t, err)

				claim2, err := VerifyFromMap(token2, conf)
				t.Log(conf.Algorithm, err)
				t.Log("map verify claim:", claim2)
				if claim["valid"].(bool) {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		}
	})

	t.Run("Verify-NonValue", func(t *testing.T) {
		for _, conf := range confArr {
			t.Log("----------------", conf.Algorithm, "----------------")
			conf.Load()
			appClaim.Issuer = ""
			token, err := Sign(appClaim, conf)
			t.Log(conf.Algorithm, err)
			t.Log("signed token:", token)
			assert.NoError(t, err)

			claim, err := Verify(token, conf)
			t.Log(conf.Algorithm, err)
			t.Log("verify claim:", claim)
			assert.NoError(t, err)

			for _, claim := range mapClaims {
				delete(claim, "iss")
				token2, err := SignWithMap(claim, conf)
				t.Log(err)
				t.Log("map signed token:", token2)
				assert.NoError(t, err)

				claim2, err := VerifyFromMap(token2, conf)
				t.Log(conf.Algorithm, err)
				t.Log("map verify claim:", claim2)
				if claim["valid"].(bool) {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			}
		}
	})
}

func TestGetToken(t *testing.T) {
	t.Run("GetToken-AppClaims", func(t *testing.T) {
		var ac = AppClaims{}
		token, err := ac.GetToken(SigningMethodRS256, "")
		assert.Equal(t, true, len(token) == 0)
		assert.Error(t, err)
	})

	t.Run("GetToken-MapClaims", func(t *testing.T) {
		var mp MapClaims = map[string]any{}
		token, err := mp.GetToken(SigningMethodRS256, "")
		assert.Equal(t, true, len(token) == 0)
		assert.Error(t, err)
	})
}

func TestSigningMethodString(t *testing.T) {
	t.Run("SigningMethodString", func(t *testing.T) {
		t.Log(SigningMethodRS256.String())
		t.Log(SigningMethodRS512.String())
		t.Log(SigningMethodHS512.String())
		t.Log(SigningMethodEdDSA.String())
		t.Log(SigningMethodES256.String())
		t.Log(SigningMethodES384.String())
		t.Log(SigningMethodES512.String())
	})
}
