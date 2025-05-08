package utils

import (
	"testing"
	"time"
)

func TestSendRequest(t *testing.T) {
	t.Run("SendRequest", func(t *testing.T) {
		_, statusCode, err := SendRequest("GET", "https://github.com", nil, nil)
		t.Log(statusCode, err)
	})

	t.Run("SendRequest-Headers", func(t *testing.T) {
		headers := map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36",
		}
		_, statusCode, err := SendRequest("GET", "https://github.com", nil, headers)
		t.Log(statusCode, err)
	})

	t.Run("SendRequest-Timeout", func(t *testing.T) {
		_, statusCode, err := SendRequest("GET", "https://github.com", nil, nil, time.Second)
		t.Log(statusCode, err)
	})

	t.Run("SendRequest-Error", func(t *testing.T) {
		_, statusCode, err := SendRequest("UNKNOWN", "https://github.com", nil, nil)
		t.Log(statusCode, err)
	})
}
