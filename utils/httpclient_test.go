package utils

import (
	"testing"
	"time"
)

func TestHttpClientImplSendRequest(t *testing.T) {
	t.Run("HttpClientImplSendRequest", func(t *testing.T) {
		_, statusCode, err := HttpClient().SendRequest("GET", "https://github.com", nil, nil)
		t.Log(statusCode, err)
	})

	t.Run("HttpClientImplSendRequest-Headers", func(t *testing.T) {
		headers := map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
		}
		_, statusCode, err := HttpClient().SendRequest("GET", "https://github.com", nil, headers)
		t.Log(statusCode, err)
	})

	t.Run("HttpClientImplSendRequest-Timeout", func(t *testing.T) {
		_, statusCode, err := HttpClient().SendRequest("GET", "https://github.com", nil, nil, time.Second)
		t.Log(statusCode, err)
	})

	t.Run("HttpClientImplSendRequest-Error", func(t *testing.T) {
		_, statusCode, err := HttpClient().SendRequest("UNKNOWN", "https://github.com", nil, nil)
		t.Log(statusCode, err)
	})
}
