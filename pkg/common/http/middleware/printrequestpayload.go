package middleware

import (
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v2/pkg/lib/logger"
	"go.uber.org/zap"
)

// PrintRequestPayload 打印请求载荷
func PrintRequestPayload() gin.HandlerFunc {
	return func(c *gin.Context) {
		dump, _ := httputil.DumpRequest(c.Request, true)
		logger.GetLogger().Info("中间件:打印请求载荷", zap.Any("value", string(dump)))

		c.Next()
	}
}
