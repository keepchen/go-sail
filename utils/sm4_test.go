package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSM4ImplECBEncrypt(t *testing.T) {
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

	for _, v := range testCases {
		result, err := SM4().ECBEncrypt(v.Key, v.Raw)
		if v.MustReturnErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, v.Result, result)
	}
}

func TestSM4ImplECBDecrypt(t *testing.T) {
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

	for _, v := range testCases {
		result, err := SM4().ECBDecrypt(v.Key, v.Raw)
		if v.MustReturnErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, v.Result, result)
	}
}
