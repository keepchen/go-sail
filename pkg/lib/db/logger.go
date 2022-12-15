package db

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type ZapLoggerForGorm struct {
	ZapLogger                 *zap.Logger
	LogLevel                  gormLogger.LogLevel
	SlowThreshold             time.Duration
	SkipCallerLookup          bool
	IgnoreRecordNotFoundError bool
}

func NewZapLoggerForGorm(zapLogger *zap.Logger, conf Conf) ZapLoggerForGorm {
	var logLevel gormLogger.LogLevel
	switch conf.LogLevel {
	case "info":
		logLevel = gormLogger.Info
	case "warn":
		logLevel = gormLogger.Warn
	case "error":
		logLevel = gormLogger.Error
	case "silent":
		logLevel = gormLogger.Silent
	default:
		logLevel = gormLogger.Info
	}
	return ZapLoggerForGorm{
		ZapLogger:                 zapLogger,
		LogLevel:                  logLevel,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: true,
	}
}

func (zg ZapLoggerForGorm) SetAsDefault() {
	gormLogger.Default = &zg
}

func (zg ZapLoggerForGorm) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := zg
	newLogger.LogLevel = level
	return &newLogger
}

func (zg ZapLoggerForGorm) Info(_ context.Context, tpl string, v ...interface{}) {
	if zg.LogLevel < gormLogger.Info {
		return
	}
	zg.ZapLogger.Sugar().Infof(tpl, v...)
}

func (zg ZapLoggerForGorm) Warn(_ context.Context, tpl string, v ...interface{}) {
	if zg.LogLevel < gormLogger.Warn {
		return
	}
	zg.ZapLogger.Sugar().Warnf(tpl, v...)
}

func (zg ZapLoggerForGorm) Error(_ context.Context, tpl string, v ...interface{}) {
	if zg.LogLevel < gormLogger.Error {
		return
	}
	zg.ZapLogger.Sugar().Errorf(tpl, v...)
}

func (zg ZapLoggerForGorm) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if zg.LogLevel < gormLogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && zg.LogLevel >= gormLogger.Error && (!zg.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		zg.ZapLogger.Error("trace", zap.Error(err), zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case zg.SlowThreshold != 0 && elapsed > zg.SlowThreshold && zg.LogLevel >= gormLogger.Warn:
		sql, rows := fc()
		zg.ZapLogger.Warn("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	case zg.LogLevel >= gormLogger.Info:
		sql, rows := fc()
		zg.ZapLogger.Debug("trace", zap.Duration("elapsed", elapsed), zap.Int64("rows", rows), zap.String("sql", sql))
	}
}
