package httpserver

import (
	"github.com/keepchen/go-sail/v3/sail/config"
	"testing"
)

func TestRunPrometheusServerWhenEnable(t *testing.T) {
	t.Run("RunPrometheusServerWhenEnable", func(t *testing.T) {
		conf := config.PrometheusConf{}
		RunPrometheusServerWhenEnable(conf)
	})

	t.Run("RunPrometheusServerWhenEnable-Enable", func(t *testing.T) {
		conf := config.PrometheusConf{
			Enable: true,
		}
		RunPrometheusServerWhenEnable(conf)
	})
}
