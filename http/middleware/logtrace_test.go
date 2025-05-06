package middleware

import (
	"net/http"
	"strings"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/logger"

	"github.com/google/uuid"
)

func TestLogTrace(t *testing.T) {
	conf := logger.Conf{
		Filename: "../../examples/logs/middleware_tester_LogTrace.log",
	}
	logger.Init(conf, "go-sail-tester")

	t.Run("LogTrace-RequestId", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("requestId", uuid.New().String())

		LogTrace()(c)
	})

	t.Run("LogTrace-XRequestId", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", uuid.New().String())

		LogTrace()(c)
	})

	t.Run("LogTrace-LengthOverflow", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")
		id := uuid.New().String()
		req.Header.Set("X-Request-ID", strings.Repeat(id, 100))

		LogTrace()(c)
	})
}
