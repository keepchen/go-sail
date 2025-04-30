package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	gormLogger "gorm.io/gorm/logger"

	"github.com/keepchen/go-sail/v3/lib/logger"
)

var (
	loggerConf = logger.Conf{
		Level:    "debug",
		Filename: "../../examples/logs/testcase_db_in_lib.log",
	}
)

func TestNewZapLoggerForGorm(t *testing.T) {
	t.Run("NewZapLoggerForGorm", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		var levels = []string{"info", "warn", "error", "silent", "unknown"}
		for _, level := range levels {
			dbConf.Logger.Level = level
			dbConf.Logger.SlowThreshold = 500
			t.Log(NewZapLoggerForGorm(logger.GetLogger(), dbConf))
		}
	})
}

func TestSetAsDefault(t *testing.T) {
	t.Run("SetAsDefault", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		zg := NewZapLoggerForGorm(logger.GetLogger(), dbConf)
		assert.NotNil(t, zg)
		zg.SetAsDefault()
		t.Log(gormLogger.Default)
	})
}

func TestLogMode(t *testing.T) {
	t.Run("LogMode", func(t *testing.T) {
		var levels = []gormLogger.LogLevel{gormLogger.Info, gormLogger.Warn, gormLogger.Error, gormLogger.Silent}

		logger.Init(loggerConf, "go-sail")
		zg := NewZapLoggerForGorm(logger.GetLogger(), dbConf)
		assert.NotNil(t, zg)

		for _, level := range levels {
			t.Log(zg.LogMode(level))
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("Info", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		zg := NewZapLoggerForGorm(logger.GetLogger(), dbConf)
		assert.NotNil(t, zg)
		zg.LogMode(gormLogger.Info)

		zg.Info(context.Background(), "tpl:%s", "var")
	})
}

func TestWarn(t *testing.T) {
	t.Run("Warn", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		zg := NewZapLoggerForGorm(logger.GetLogger(), dbConf)
		assert.NotNil(t, zg)
		zg.LogMode(gormLogger.Warn)

		zg.Warn(context.Background(), "tpl:%s", "var")
	})
}

func TestError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		zg := NewZapLoggerForGorm(logger.GetLogger(), dbConf)
		assert.NotNil(t, zg)
		zg.LogMode(gormLogger.Error)

		zg.Error(context.Background(), "tpl:%s", "var")
	})
}

func TestTrace(t *testing.T) {
	t.Run("Trace", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		zg := NewZapLoggerForGorm(logger.GetLogger(), dbConf)
		assert.NotNil(t, zg)
		var (
			levels = []gormLogger.LogLevel{gormLogger.Info, gormLogger.Warn, gormLogger.Error, gormLogger.Silent}
			fn     = func() (sql string, rowsAffected int64) {
				return "select * from user", 1
			}
		)

		for _, level := range levels {
			t.Log(zg.LogMode(level))
			zg.Trace(context.Background(), time.Now(), fn, nil)
		}

	})
}
