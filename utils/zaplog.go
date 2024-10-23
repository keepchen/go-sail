package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

// LogDebug 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：debug
func LogDebug(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).Debug(message, fields...)
		return
	}
	logger.GetLogger().Debug(message, fields...)
}

// LogInfo 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：info
func LogInfo(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).Info(message, fields...)
		return
	}
	logger.GetLogger().Info(message, fields...)
}

// LogWarn 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：warn
func LogWarn(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).Warn(message, fields...)
		return
	}
	logger.GetLogger().Warn(message, fields...)
}

// LogError 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：error
func LogError(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).Error(message, fields...)
		return
	}
	logger.GetLogger().Error(message, fields...)
}

// LogDPanic 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：dpanic
func LogDPanic(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).DPanic(message, fields...)
		return
	}
	logger.GetLogger().DPanic(message, fields...)
}

// LogPanic 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：panic
func LogPanic(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).Panic(message, fields...)
		return
	}
	logger.GetLogger().Panic(message, fields...)
}

// LogFatal 打印日志
//
// 尝试从上下文中获取logger实例进行打印，
//
// 从上下文中获取的logger实例会携带链路追踪的字段。
//
// 若获取失败，则调用logger库进行打印
//
// 日志级别：fatal
func LogFatal(c *gin.Context, message string, fields ...zap.Field) {
	if loggerSvc, ok := c.Get("logger"); ok {
		loggerSvc.(*zap.Logger).Fatal(message, fields...)
		return
	}
	logger.GetLogger().Fatal(message, fields...)
}
