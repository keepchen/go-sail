package sail

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/stretchr/testify/require"
)

func TestSailRSAOAEP(t *testing.T) {
	previousConfig := config.Get()
	t.Cleanup(func() { config.Set(previousConfig) })
	config.Set(&config.Config{JwtConf: &jwt.Conf{
		PublicKey:  string(publicKey),
		PrivateKey: string(privateKey),
	}})

	keyBlock, _ := pem.Decode(privateKey)
	require.NotNil(t, keyBlock)
	key, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	require.NoError(t, err)
	priv, ok := key.(*rsa.PrivateKey)
	require.True(t, ok)
	plaintext := "sail OAEP interoperability"
	standardCiphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &priv.PublicKey, []byte(plaintext), nil)
	require.NoError(t, err)

	for _, tc := range []struct {
		name    string
		encrypt func(string) (string, error)
		decrypt func(string) (string, error)
	}{
		{
			name: "Password",
			encrypt: func(s string) (string, error) {
				return Password().RSAEncryptOAEP(s, publicKey)
			},
			decrypt: func(s string) (string, error) {
				return Password().RSADecryptOAEP(s, privateKey)
			},
		},
		{name: "JWT", encrypt: JWT().EncryptOAEP, decrypt: JWT().DecryptOAEP},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ciphertext, err := tc.encrypt(plaintext)
			require.NoError(t, err)
			decoded, err := base64.StdEncoding.DecodeString(ciphertext)
			require.NoError(t, err)
			standardPlaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, decoded, nil)
			require.NoError(t, err)
			require.Equal(t, plaintext, string(standardPlaintext))

			result, err := tc.decrypt(ciphertext)
			require.NoError(t, err)
			require.Equal(t, plaintext, result)
			result, err = tc.decrypt(base64.StdEncoding.EncodeToString(standardCiphertext))
			require.NoError(t, err)
			require.Equal(t, plaintext, result)

			result, err = tc.encrypt(strings.Repeat("a", 191))
			require.ErrorIs(t, err, rsa.ErrMessageTooLong)
			require.Empty(t, result)
			result, err = tc.decrypt("!invalid base64!")
			require.Error(t, err)
			require.Empty(t, result)
		})
	}
}

func TestSailRSAOAEPInvalidKeys(t *testing.T) {
	previousConfig := config.Get()
	t.Cleanup(func() { config.Set(previousConfig) })
	config.Set(&config.Config{JwtConf: &jwt.Conf{
		PublicKey:  "invalid public key",
		PrivateKey: "invalid private key",
	}})

	result, err := JWT().EncryptOAEP("message")
	require.Error(t, err)
	require.Empty(t, result)
	result, err = JWT().DecryptOAEP(base64.StdEncoding.EncodeToString([]byte("ciphertext")))
	require.Error(t, err)
	require.Empty(t, result)
	result, err = Password().RSAEncryptOAEP("message", nil)
	require.Error(t, err)
	require.Empty(t, result)
	result, err = Password().RSADecryptOAEP(base64.StdEncoding.EncodeToString([]byte("ciphertext")), nil)
	require.Error(t, err)
	require.Empty(t, result)
}
