package db

import (
	"testing"

	"github.com/keepchen/go-sail/v3/lib/logger"
)

var (
	loggerConf = logger.Conf{
		Level:    "debug",
		Filename: "../../examples/logs/testcase_db_in_lib.log",
	}
)

func TestNewZapLoggerForGorm(t *testing.T) {
	logger.Init(loggerConf, "go-sail")
	t.Run("NewZapLoggerForGorm", func(t *testing.T) {
		t.Log(NewZapLoggerForGorm(logger.GetLogger(), dbConf))
	})
}
