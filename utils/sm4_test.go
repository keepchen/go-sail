package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSM4ImplEncrypt(t *testing.T) {
	testCases := []struct {
		Key           string
		Raw           string
		Result        string
		MustReturnErr bool
	}{
		{
			Key:           "1233",
			Raw:           "12345678",
			MustReturnErr: true,
		},
		{
			Key:    "c6ee9fe905390c7dc30f124af5b18771",
			Raw:    "123456",
			Result: "aQomDVTfAKDhKgwED8a1ug==",
		},
		{
			Key:    "c6ee9fe905390c7dc30f124af5b18771",
			Raw:    "hello,world!你好，世界！",
			Result: "RJRzbYKvf6zsY7Bx0QaMjvHl5Y3y/Xj0CIpG8bED5b4=",
		},
	}

	t.Run("ECBEncrypt", func(t *testing.T) {
		for _, v := range testCases {
			result, err := SM4().ECBEncrypt(v.Key, v.Raw)
			if v.MustReturnErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v.Result, result)
		}
	})

	t.Run("ECBEncrypt-Error", func(t *testing.T) {
		for _, v := range testCases {
			_, err := SM4().ECBEncrypt(v.Key+"?", v.Raw)
			assert.Error(t, err)
			//assert.Equal(t, v.Result, result)
		}
	})

	t.Run("GCMEncrypt", func(t *testing.T) {
		for _, v := range testCases {
			key := v.Key
			if !v.MustReturnErr {
				key = v.Key[:16]
			}
			result, err := SM4().GCMEncrypt(key, v.Raw)
			if v.MustReturnErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		}
	})
}

func TestSM4ImplDecrypt(t *testing.T) {
	testCases := []struct {
		Key           string
		Raw           string
		Result        string
		MustReturnErr bool
	}{
		{
			Key:           "1233",
			Raw:           "12345678",
			MustReturnErr: true,
		},
		{
			Key:    "c6ee9fe905390c7dc30f124af5b18771",
			Raw:    "aQomDVTfAKDhKgwED8a1ug==",
			Result: "123456",
		},
		{
			Key:    "c6ee9fe905390c7dc30f124af5b18771",
			Raw:    "RJRzbYKvf6zsY7Bx0QaMjvHl5Y3y/Xj0CIpG8bED5b4=",
			Result: "hello,world!你好，世界！",
		},
	}

	t.Run("ECBDecrypt", func(t *testing.T) {
		for _, v := range testCases {
			result, err := SM4().ECBDecrypt(v.Key, v.Raw)
			if v.MustReturnErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, v.Result, result)
		}
	})

	t.Run("ECBDecrypt-Error", func(t *testing.T) {
		for _, v := range testCases {
			_, err := SM4().ECBDecrypt(v.Key+"?", v.Raw)
			assert.Error(t, err)
			//assert.NotEqual(t, v.Result, result)
		}
	})

	t.Run("GCMDecrypt", func(t *testing.T) {
		for _, v := range testCases {
			key := v.Key
			if !v.MustReturnErr {
				key = v.Key[:16]
			}
			result, err1 := SM4().GCMEncrypt(key, v.Raw)
			r2, err2 := SM4().GCMDecrypt(key, result)
			if v.MustReturnErr {
				assert.Error(t, err1)
				assert.Empty(t, result)
				assert.Error(t, err2)
				assert.Empty(t, r2)
			} else {
				assert.NoError(t, err1)
				assert.NotEmpty(t, result)
				assert.NoError(t, err2)
				assert.NotEmpty(t, r2)
			}
		}
	})
}
