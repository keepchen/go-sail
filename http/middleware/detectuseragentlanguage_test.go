package middleware

import (
	"net/http"
	"testing"
)

func TestDetectUserAgentLanguage(t *testing.T) {
	t.Run("DetectUserAgentLanguage", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7")

		DetectUserAgentLanguage()(c)
	})
}
