package sail

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

// LogTracer 链路日志追踪器
type LogTracer interface {
	// GetLogger 获取 zap.Logger 实例
	//
	// 此实例携带了上下文的requestId和spanId，用于链路追踪
	GetLogger() *zap.Logger
	// Debug 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：Debug
	Debug(message string, fields ...zap.Field)
	// Info 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：Info
	Info(message string, fields ...zap.Field)
	// Warn 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：Warn
	Warn(message string, fields ...zap.Field)
	// Error 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：Error
	Error(message string, fields ...zap.Field)
	// DPanic 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：DPanic
	DPanic(message string, fields ...zap.Field)
	// Panic 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：Panic
	Panic(message string, fields ...zap.Field)
	// Fatal 打印日志
	//
	// 尝试从上下文中获取logger实例进行打印，
	//
	// 从上下文中获取的logger实例会携带链路追踪的字段。
	//
	// 若获取失败，则调用logger库进行打印
	//
	// 日志级别：Fatal
	Fatal(message string, fields ...zap.Field)
}

type logTrace struct {
	ginContext *gin.Context
	loggerSvc  *zap.Logger
}

var _ LogTracer = (*logTrace)(nil)

// LogTrace 链路日志追踪
func LogTrace(c *gin.Context) LogTracer {
	var zapLogger *zap.Logger
	if loggerSvc, ok := c.Get("logger"); ok {
		zapLogger = loggerSvc.(*zap.Logger)
	} else {
		zapLogger = logger.GetLogger()
	}

	return &logTrace{ginContext: c, loggerSvc: zapLogger}
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
