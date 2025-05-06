package middleware

import (
	"net/http"
	"testing"
	"time"

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

		c.Request = req

		PrometheusExporter()(c)
	})

	t.Run("PrometheusExporter-Write", func(t *testing.T) {
		c, r := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		c.Request = req

		r.Use(PrometheusExporter())

		c.String(http.StatusOK, "%s", "OK")
	})
}

func TestSetDiskPath(t *testing.T) {
	t.Run("SetDiskPath", func(t *testing.T) {
		SetDiskPath("/data")
	})
}

func TestSetSampleInterval(t *testing.T) {
	t.Run("SetSampleInterval-OK", func(t *testing.T) {
		SetSampleInterval("10s")
	})

	t.Run("SetSampleInterval-Error", func(t *testing.T) {
		SetSampleInterval("10xx")
	})
}

func TestSystemMetricsSample(t *testing.T) {
	t.Run("SystemMetricsSample", func(t *testing.T) {
		SetSampleInterval("5s")
		SetDiskPath("/")
		exitChan := make(chan struct{})
		go func() {
			time.Sleep(time.Second * 10)
			exitChan <- struct{}{}
		}()
		SystemMetricsSample(exitChan)
	})
}
