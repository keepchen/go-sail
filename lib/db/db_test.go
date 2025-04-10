package db

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/lib/logger"

	"github.com/keepchen/go-sail/v3/constants"
)

var (
	dbConf = Conf{
		Enable:      true,
		DriverName:  "mysql",
		AutoMigrate: true,
		Logger: Logger{
			Level: "debug",
		},
		ConnectionPool: ConnectionPoolConf{
			MaxOpenConnCount:       100,
			MaxIdleConnCount:       10,
			ConnMaxLifeTimeMinutes: 30,
			ConnMaxIdleTimeMinutes: 10,
		},
		Mysql: MysqlConf{
			Read: MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      33060,
				Username:  "root",
				Password:  "changeMe",
				Database:  "go-sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
			Write: MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      33060,
				Username:  "root",
				Password:  "changeMe",
				Database:  "go-sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
		},
	}
)

func TestInitDB(t *testing.T) {
	logger.Init(loggerConf, "go-sail")
	dbr, _, dbw, _ := New(dbConf)
	if dbr == nil || dbw == nil {
		t.Log("database instance is nil, testing not emit.")
		return
	}
	t.Run("InitDB", func(t *testing.T) {
		InitDB(dbConf)
	})
}

func TestGetInstance(t *testing.T) {
	assert.Nil(t, GetInstance())
}
