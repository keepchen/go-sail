package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRSAOAEPRoundTrip(t *testing.T) {
	for _, keys := range []struct {
		name    string
		public  []byte
		private []byte
	}{
		{name: "PEM", public: publicKey, private: privateKey},
		{name: "inline", public: publicKeyInline, private: privateKeyInline},
	} {
		t.Run(keys.name, func(t *testing.T) {
			for _, plaintext := range []string{"", "password", "你好，OAEP", strings.Repeat("a", 190)} {
				ciphertext, err := RSA().EncryptOAEP(plaintext, keys.public, nil)
				require.NoError(t, err)
				decoded, err := base64.StdEncoding.DecodeString(ciphertext)
				require.NoError(t, err)
				require.Len(t, decoded, 256)
				result, err := RSA().DecryptOAEP(ciphertext, keys.private, nil)
				require.NoError(t, err)
				require.Equal(t, plaintext, result)
			}
		})
	}
}

func TestRSAOAEPStandardLibraryCompatibility(t *testing.T) {
	pub := parsePublicKey(publicKey)
	priv := parsePrivateKey(privateKey)
	require.NotNil(t, pub)
	require.NotNil(t, priv)
	plaintext := "OAEP interoperability"
	label := []byte("go-sail:interoperability")

	ciphertext, err := RSA().EncryptOAEP(plaintext, publicKey, label)
	require.NoError(t, err)
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)
	result, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, decoded, label)
	require.NoError(t, err)
	require.Equal(t, plaintext, string(result))

	encoded, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, []byte(plaintext), label)
	require.NoError(t, err)
	resultString, err := RSA().DecryptOAEP(base64.StdEncoding.EncodeToString(encoded), privateKey, label)
	require.NoError(t, err)
	require.Equal(t, plaintext, resultString)

	otherCiphertext, err := RSA().EncryptOAEP(plaintext, publicKey, label)
	require.NoError(t, err)
	require.NotEqual(t, ciphertext, otherCiphertext)
}

func TestRSAOAEPLabelBinding(t *testing.T) {
	label := []byte("go-sail:password")
	ciphertext, err := RSA().EncryptOAEP("message", publicKey, label)
	require.NoError(t, err)
	result, err := RSA().DecryptOAEP(ciphertext, privateKey, label)
	require.NoError(t, err)
	require.Equal(t, "message", result)

	for _, wrongLabel := range [][]byte{nil, {}, []byte("go-sail:token")} {
		result, err := RSA().DecryptOAEP(ciphertext, privateKey, wrongLabel)
		require.ErrorIs(t, err, rsa.ErrDecryption)
		require.Empty(t, result)
	}

	for _, emptyLabel := range [][]byte{nil, {}} {
		ciphertext, err := RSA().EncryptOAEP("message", publicKey, emptyLabel)
		require.NoError(t, err)
		for _, equivalentLabel := range [][]byte{nil, {}} {
			result, err := RSA().DecryptOAEP(ciphertext, privateKey, equivalentLabel)
			require.NoError(t, err)
			require.Equal(t, "message", result)
		}
	}
}

func TestRSAOAEPRejectsInvalidInput(t *testing.T) {
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	ecPublicDER, err := x509.MarshalPKIXPublicKey(&ecKey.PublicKey)
	require.NoError(t, err)
	ecPrivateDER, err := x509.MarshalPKCS8PrivateKey(ecKey)
	require.NoError(t, err)

	for _, tc := range []struct {
		name string
		key  []byte
	}{
		{name: "empty"},
		{name: "malformed", key: []byte("invalid key")},
		{name: "non-RSA", key: pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: ecPublicDER})},
	} {
		t.Run("public/"+tc.name, func(t *testing.T) {
			result, err := RSA().EncryptOAEP("message", tc.key, nil)
			require.Error(t, err)
			require.Empty(t, result)
		})
	}

	ciphertext, err := RSA().EncryptOAEP("message", publicKey, nil)
	require.NoError(t, err)
	for _, tc := range []struct {
		name string
		key  []byte
	}{
		{name: "empty"},
		{name: "malformed", key: []byte("invalid key")},
		{name: "non-RSA", key: pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: ecPrivateDER})},
	} {
		t.Run("private/"+tc.name, func(t *testing.T) {
			result, err := RSA().DecryptOAEP(ciphertext, tc.key, nil)
			require.Error(t, err)
			require.Empty(t, result)
		})
	}

	t.Run("too long", func(t *testing.T) {
		result, err := RSA().EncryptOAEP(strings.Repeat("a", 191), publicKey, nil)
		require.ErrorIs(t, err, rsa.ErrMessageTooLong)
		require.Empty(t, result)
	})

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)
	decoded[len(decoded)-1] ^= 1
	for _, invalid := range []string{"", "!invalid base64!", base64.StdEncoding.EncodeToString([]byte("short")), base64.StdEncoding.EncodeToString(decoded)} {
		result, err := RSA().DecryptOAEP(invalid, privateKey, nil)
		require.Error(t, err)
		require.Empty(t, result)
	}

	t.Run("wrong key", func(t *testing.T) {
		otherKey, err := rsa.GenerateKey(rand.Reader, 2048)
		require.NoError(t, err)
		keyDER, err := x509.MarshalPKCS8PrivateKey(otherKey)
		require.NoError(t, err)
		keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
		result, err := RSA().DecryptOAEP(ciphertext, keyPEM, nil)
		require.ErrorIs(t, err, rsa.ErrDecryption)
		require.Empty(t, result)
	})
}

func TestRSAOAEPRejectsDifferentParameters(t *testing.T) {
	pub := parsePublicKey(publicKey)
	require.NotNil(t, pub)
	wrongHash, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pub, []byte("message"), nil)
	require.NoError(t, err)
	wrongLabel, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, []byte("message"), []byte("other purpose"))
	require.NoError(t, err)
	for _, ciphertext := range [][]byte{wrongHash, wrongLabel} {
		result, err := RSA().DecryptOAEP(base64.StdEncoding.EncodeToString(ciphertext), privateKey, nil)
		require.ErrorIs(t, err, rsa.ErrDecryption)
		require.Empty(t, result)
	}
}

func TestRSALegacyCompatibility(t *testing.T) {
	pub := parsePublicKey(publicKey)
	priv := parsePrivateKey(privateKey)
	require.NotNil(t, pub)
	require.NotNil(t, priv)
	plaintext := "legacy ciphertext"
	ciphertext, err := RSA().Encrypt(plaintext, publicKey)
	require.NoError(t, err)
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	require.NoError(t, err)
	result, err := rsa.DecryptPKCS1v15(rand.Reader, priv, decoded) //nolint:staticcheck // Verify the deprecated API still produces legacy-compatible ciphertext.
	require.NoError(t, err)
	require.Equal(t, plaintext, string(result))

	encoded, err := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(plaintext)) //nolint:staticcheck // Generate legacy ciphertext independently for compatibility testing.
	require.NoError(t, err)
	resultString, err := RSA().Decrypt(base64.StdEncoding.EncodeToString(encoded), privateKey)
	require.NoError(t, err)
	require.Equal(t, plaintext, resultString)

	resultString, err = RSA().DecryptOAEP(ciphertext, privateKey, nil)
	require.ErrorIs(t, err, rsa.ErrDecryption)
	require.Empty(t, resultString)
	oaepCiphertext, err := RSA().EncryptOAEP(plaintext, publicKey, nil)
	require.NoError(t, err)
	resultString, err = RSA().Decrypt(oaepCiphertext, privateKey)
	require.ErrorIs(t, err, rsa.ErrDecryption)
	require.Empty(t, resultString)
}
