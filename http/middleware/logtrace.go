package middleware

import (
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

// LogTrace 日志追踪
//
// 作用是在请求入口注入必要的内容到上下文，供后续的请求调用链使用，一般用于日志追踪、链路追踪
//
// 注入内容：
//
// 1.请求id
//
// 2.请求到来那一刻的纳秒时间戳
//
// 3.包装了请求id的日志组件实例
//
// Example:
//
// requestId, ok := ginContext.Get("requestId").(string)
//
// spanId, ok := ginContext.Get("spanId").(string)
//
// entryAt, ok := ginContext.Get("entryAt").(int64)
//
// logger, ok := ginContext.Get("logger").(*zap.Logger)
func LogTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			requestId string
			spanId    string
		)
		requestIdInHeader := c.Request.Header.Get("requestId")
		if len(requestIdInHeader) == 0 {
			requestIdInHeader = c.Request.Header.Get("X-Request-Id")
		}

		//溢出检测
		//TODO 目前暂定为最长保留200个字符
		if utf8.RuneCountInString(requestIdInHeader) > 200 {
			requestIdInHeader = string([]rune(requestIdInHeader)[:200])
		}

		if len(requestIdInHeader) > 0 {
			requestId = requestIdInHeader
			spanId = uuid.New().String()
		} else {
			requestId = uuid.New().String()
			spanId = requestId
		}
		c.Set("requestId", requestId)
		c.Set("spanId", spanId)
		c.Set("entryAt", time.Now().UnixNano())
		c.Set("logger", logger.GetLogger().With(zap.String("requestId", requestId),
			zap.String("spanId", spanId)))

		c.Next()
	}
}
