package sail

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/sail/config"
)

func TestPrintSummaryInfo(t *testing.T) {
	t.Run("printSummaryInfo", func(t *testing.T) {
		conf := config.Config{}
		ginEngine := gin.Default()
		printSummaryInfo(conf, ginEngine)
	})

	t.Run("printSummaryInfo-Debug", func(t *testing.T) {
		conf := config.Config{
			HttpServer: config.HttpServerConf{
				Debug: true,
				Prometheus: config.PrometheusConf{
					Enable: true,
				},
				Swagger: config.SwaggerConf{
					Enable: true,
				},
			},
		}
		ginEngine := gin.Default()
		printSummaryInfo(conf, ginEngine)
	})
}
