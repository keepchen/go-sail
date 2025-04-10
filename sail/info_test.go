package sail

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/sail/config"
	"testing"
)

func TestPrintSummaryInfo(t *testing.T) {
	t.Run("printSummaryInfo", func(t *testing.T) {
		conf := config.Config{}
		ginEngine := gin.Default()
		printSummaryInfo(conf.HttpServer, ginEngine)
	})
}
