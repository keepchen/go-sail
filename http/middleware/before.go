package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

// Before 中间件拦截器-before
func Before() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := uuid.New().String()
		c.Set("requestId", requestId)
		c.Set("logger", logger.GetLogger().With(zap.String("requestId", requestId)))

		c.Next()
	}
}
