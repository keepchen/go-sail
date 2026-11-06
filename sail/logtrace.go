package sail

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

// LogTracer log tracer for call chains
type LogTracer interface {
	// RequestID retrieves the request ID
	//
	// Try to retrieve the request ID from the context
	RequestID() string
	// SpanID retrieves the internal span ID
	//
	// Try to retrieve the internal span ID from the context
	SpanID() string
	// EntryAt retrieves the request arrival timestamp in nanoseconds
	//
	// Try to retrieve the request arrival timestamp in nanoseconds from the context
	EntryAt() int64
	// GetLogger retrieves the zap.Logger instance
	//
	// This instance includes the context's requestId and spanId for tracing
	GetLogger() *zap.Logger
	// Debug prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: Debug
	Debug(message string, fields ...zap.Field)
	// Info prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: Info
	Info(message string, fields ...zap.Field)
	// Warn prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: Warn
	Warn(message string, fields ...zap.Field)
	// Error prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: Error
	Error(message string, fields ...zap.Field)
	// DPanic prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: DPanic
	DPanic(message string, fields ...zap.Field)
	// Panic prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: Panic
	Panic(message string, fields ...zap.Field)
	// Fatal prints a log
	//
	// Try to retrieve a logger instance from the context for logging;
	//
	// the context logger includes tracing fields.
	//
	// If retrieval fails, use the logger library to print
	//
	// Log level: Fatal
	Fatal(message string, fields ...zap.Field)
}

type logTrace struct {
	ginContext *gin.Context
	loggerSvc  *zap.Logger
}

var _ LogTracer = (*logTrace)(nil)

// LogTrace call chain log tracing
func LogTrace(c *gin.Context) LogTracer {
	var zapLogger *zap.Logger
	if loggerSvc, ok := c.Get("logger"); ok {
		zapLogger = loggerSvc.(*zap.Logger)
	} else {
		zapLogger = logger.GetLogger()
	}

	return &logTrace{ginContext: c, loggerSvc: zapLogger}
}

func (l *logTrace) RequestID() string {
	var requestID string
	if val, ok := l.ginContext.Get("requestId"); ok {
		requestID, _ = val.(string)
	}

	return requestID
}

func (l *logTrace) SpanID() string {
	var spanID string
	if val, ok := l.ginContext.Get("spanId"); ok {
		spanID, _ = val.(string)
	}

	return spanID
}

func (l *logTrace) EntryAt() int64 {
	var entryAt int64
	if val, ok := l.ginContext.Get("entryAt"); ok {
		entryAt, _ = val.(int64)
	}

	return entryAt
}

func (l *logTrace) GetLogger() *zap.Logger {
	return l.loggerSvc
}

func (l *logTrace) Debug(message string, fields ...zap.Field) {
	l.loggerSvc.Debug(message, fields...)
}

func (l *logTrace) Info(message string, fields ...zap.Field) {
	l.loggerSvc.Info(message, fields...)
}

func (l *logTrace) Warn(message string, fields ...zap.Field) {
	l.loggerSvc.Warn(message, fields...)
}

func (l *logTrace) Error(message string, fields ...zap.Field) {
	l.loggerSvc.Error(message, fields...)
}

func (l *logTrace) DPanic(message string, fields ...zap.Field) {
	l.loggerSvc.DPanic(message, fields...)
}

func (l *logTrace) Panic(message string, fields ...zap.Field) {
	l.loggerSvc.Panic(message, fields...)
}

func (l *logTrace) Fatal(message string, fields ...zap.Field) {
	l.loggerSvc.Fatal(message, fields...)
}
