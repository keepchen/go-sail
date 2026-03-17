package sail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	plaintext = "password"
	hexKey    = "c6ee9fe905390c7dc30f124af5b18771"
)

func TestPassword(t *testing.T) {
	t.Run("Password", func(t *testing.T) {
		assert.NotNil(t, Password())
	})
}

func TestAesGCMEncrypt(t *testing.T) {
	t.Run("AesGCMEncrypt", func(t *testing.T) {
		result, err := Password().AesGCMEncrypt(plaintext, hexKey)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestAesGCMDecrypt(t *testing.T) {
	t.Run("AesGCMDecrypt", func(t *testing.T) {
		encryptedStr, err := Password().AesGCMEncrypt(plaintext, hexKey)
		assert.Nil(t, err)
		assert.NotEmpty(t, encryptedStr)
		decrypted, err := Password().AesGCMDecrypt(encryptedStr, hexKey)
		assert.Nil(t, err)
		assert.Equal(t, plaintext, decrypted)
	})
}

func TestSM4GCMEncrypt(t *testing.T) {
	t.Run("SM4GCMEncrypt", func(t *testing.T) {
		result, err := Password().SM4GCMEncrypt(plaintext, hexKey[:16])
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestSM4GCMDecrypt(t *testing.T) {
	t.Run("SM4GCMDecrypt", func(t *testing.T) {
		encryptedStr, err := Password().SM4GCMEncrypt(plaintext, hexKey[:16])
		assert.Nil(t, err)
		assert.NotEmpty(t, encryptedStr)
		decrypted, err := Password().SM4GCMDecrypt(encryptedStr, hexKey[:16])
		assert.Nil(t, err)
		assert.Equal(t, plaintext, decrypted)
	})
}

func TestRSAEncrypt(t *testing.T) {
	t.Run("RSAEncrypt", func(t *testing.T) {
		result, err := Password().RSAEncrypt(plaintext, publicKey)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestRSADecrypt(t *testing.T) {
	t.Run("RSADecrypt", func(t *testing.T) {
		result, err := Password().RSAEncrypt(plaintext, publicKey)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)

		result, err = Password().RSADecrypt(result, privateKey)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestBcryptHashMake(t *testing.T) {
	t.Run("BcryptHashMake", func(t *testing.T) {
		result, err := Password().BcryptHashMake(plaintext)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestBcryptHashCheck(t *testing.T) {
	t.Run("BcryptHashCheck", func(t *testing.T) {
		result, err := Password().BcryptHashMake(plaintext)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)

		valid, err := Password().BcryptHashCheck(plaintext, result)
		assert.Nil(t, err)
		assert.True(t, valid)
	})
}

func TestArgon2HashMake(t *testing.T) {
	t.Run("Argon2HashMake", func(t *testing.T) {
		result, err := Password().Argon2HashMake(plaintext)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
	})
}

func TestArgon2HashCheck(t *testing.T) {
	t.Run("Argon2HashCheck", func(t *testing.T) {
		result, err := Password().Argon2HashMake(plaintext)
		assert.Nil(t, err)
		assert.NotEmpty(t, result)
		valid, err := Password().Argon2HashCheck(plaintext, result)
		assert.Nil(t, err)
		assert.True(t, valid)
	})
}
