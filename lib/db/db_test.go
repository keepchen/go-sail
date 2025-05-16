package db

import (
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

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
		NowFunc: func() time.Time {
			return time.Now().In(time.UTC)
		},
		Mysql: MysqlConf{
			Read: MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      3306,
				Username:  "root",
				Password:  "root",
				Database:  "go_sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
			Write: MysqlConfItem{
				Host:      "127.0.0.1",
				Port:      3306,
				Username:  "root",
				Password:  "root",
				Database:  "go_sail",
				Charset:   "utf8mb4",
				ParseTime: true,
				Loc:       url.QueryEscape(constants.TimeZoneUTCPlus7),
			},
		},
	}
)

func TestNewFreshDB(t *testing.T) {
	t.Run("NewFreshDB", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, err1, dbw, err2 := NewFreshDB(dbConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotNil(t, dbr)
		assert.NotNil(t, dbw)
		r, _ := dbr.DB()
		_ = r.Close()
		w, _ := dbw.DB()
		_ = w.Close()
	})
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		dbr, err1, dbw, err2 := New(dbConf)
		if dbr == nil || dbw == nil {
			t.Log("database instance is nil, testing not emit.")
			return
		}
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotNil(t, dbr)
		assert.NotNil(t, dbw)
		r, _ := dbr.DB()
		_ = r.Close()
		w, _ := dbw.DB()
		_ = w.Close()
	})
}

func TestInitDB(t *testing.T) {
	t.Run("InitDB-Panic", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")

		//unknown port
		dbConf.Mysql.Read.Port += 1
		dbConf.Mysql.Write.Port += 1

		assert.Panics(t, func() {
			InitDB(dbConf)
		})
	})

	t.Run("InitDB", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		InitDB(dbConf)
		assert.NotNil(t, GetInstance())
		r, _ := GetInstance().W.DB()
		_ = r.Close()
		w, _ := GetInstance().W.DB()
		_ = w.Close()
	})
}

func TestInit(t *testing.T) {
	t.Run("Init-Panic", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")

		//unknown port
		dbConf.Mysql.Read.Port += 1
		dbConf.Mysql.Write.Port += 1

		assert.Panics(t, func() {
			Init(dbConf)
			assert.NotNil(t, GetInstance())
			r, _ := GetInstance().W.DB()
			_ = r.Close()
			w, _ := GetInstance().W.DB()
			_ = w.Close()
		})
	})

	t.Run("Init", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(dbConf)
		assert.NotNil(t, GetInstance())
		r, _ := GetInstance().W.DB()
		_ = r.Close()
		w, _ := GetInstance().W.DB()
		_ = w.Close()
	})
}

func TestMustInitDB(t *testing.T) {
	t.Run("mustInitDB-Panic", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")

		//unknown port
		dbConf.Mysql.Read.Port += 1
		dbConf.Mysql.Write.Port += 1

		assert.Panics(t, func() {
			d1, d2 := dbConf.GenDialector()
			r := mustInitDB(dbConf, d1)
			w := mustInitDB(dbConf, d2)
			assert.NotNil(t, r)
			assert.NotNil(t, w)
			rd, _ := r.DB()
			_ = rd.Close()
			wd, _ := w.DB()
			_ = wd.Close()
		})
	})

	t.Run("mustInitDB", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		d1, d2 := dbConf.GenDialector()
		r := mustInitDB(dbConf, d1)
		w := mustInitDB(dbConf, d2)
		assert.NotNil(t, r)
		assert.NotNil(t, w)
		rd, _ := r.DB()
		_ = rd.Close()
		wd, _ := w.DB()
		_ = wd.Close()
	})
}

func TestOnlyInitDB(t *testing.T) {
	t.Run("initDB-Failure", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")

		//unknown port
		dbConf.Mysql.Read.Port += 1
		dbConf.Mysql.Write.Port += 1

		d1, d2 := dbConf.GenDialector()
		r, e1 := initDB(dbConf, d1)
		w, e2 := initDB(dbConf, d2)
		assert.Error(t, e1)
		assert.Error(t, e2)
		assert.Nil(t, r)
		assert.Nil(t, w)
	})

	t.Run("initDB", func(t *testing.T) {
		logger.Init(loggerConf, "go-sail")
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		d1, d2 := dbConf.GenDialector()
		r, e1 := initDB(dbConf, d1)
		w, e2 := initDB(dbConf, d2)
		assert.NoError(t, e1)
		assert.NoError(t, e2)
		assert.NotNil(t, r)
		assert.NotNil(t, w)
		rd, _ := r.DB()
		_ = rd.Close()
		wd, _ := w.DB()
		_ = wd.Close()
	})
}

func TestGetInstance(t *testing.T) {
	assert.Nil(t, GetInstance())
}
