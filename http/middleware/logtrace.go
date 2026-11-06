package middleware

import (
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

type tracerStd struct {
	requestId string
	spanId    string
}

var tracerStdPool = sync.Pool{
	New: func() interface{} {
		return &tracerStd{}
	},
}

// LogTrace log tracing
//
// Inject necessary values into the context at the request entry point for subsequent calls, typically for log and distributed tracing
//
// Injected values:
//
// 1.Request ID
//
// 2.Nanosecond timestamp when the request arrives
//
// 3.Logger instance containing the request ID
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
		tracer := tracerStdPool.Get().(*tracerStd)
		tracer.requestId = ""
		tracer.spanId = ""

		tracer.requestId = c.Request.Header.Get("requestId")
		if len(tracer.requestId) == 0 {
			tracer.requestId = c.Request.Header.Get("X-Request-Id")
		}

		//Overflow check
		//TODO Currently limited to a maximum of 200 characters
		if utf8.RuneCountInString(tracer.requestId) > 200 {
			tracer.requestId = string([]rune(tracer.requestId)[:200])
		}

		if len(tracer.requestId) > 0 {
			tracer.spanId = uuid.New().String()
		} else {
			tracer.requestId = uuid.New().String()
			tracer.spanId = tracer.requestId
		}
		c.Set("requestId", tracer.requestId)
		c.Set("spanId", tracer.spanId)
		c.Set("entryAt", time.Now().UnixNano())
		c.Set("logger", logger.GetLogger().With(zap.String("requestId", tracer.requestId),
			zap.String("spanId", tracer.spanId)))

		//release
		tracerStdPool.Put(tracer)

		c.Next()
	}
}
