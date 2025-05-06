package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 创建gin测试下上文和引擎
func createTestContextAndEngine() (*gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()

	//创建测试用的Request（可自定义请求参数）
	req, _ := http.NewRequest("GET", "/test?name=foo", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", uuid.New().String())
	req.Header.Set("requestId", uuid.New().String())
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7,la;q=0.6")
	c, r := gin.CreateTestContext(w)
	c.Request = req
	c.Set("requestId", uuid.New().String())

	return c, r
}

func TestDetectRealClientIP(t *testing.T) {
	t.Run("DetectRealClientIP-XFF", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "127.0.0.1")

		c.Request = req

		DetectRealClientIP()(c)
	})

	t.Run("DetectRealClientIP-XFF-None", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "")

		c.Request = req

		DetectRealClientIP()(c)
	})

	t.Run("DetectRealClientIP-XFF-Malformed", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Forwarded-For", "127.0.0.1.")

		c.Request = req

		DetectRealClientIP()(c)
	})
}
