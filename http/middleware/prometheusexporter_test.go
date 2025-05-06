package middleware

import (
	"net/http"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/logger"
)

func TestPrometheusExporter(t *testing.T) {
	conf := logger.Conf{
		Filename: "../../examples/logs/middleware_tester_PrometheusExporter.log",
	}
	logger.Init(conf, "go-sail-tester")

	t.Run("PrometheusExporter", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		PrometheusExporter()(c)
	})
}
