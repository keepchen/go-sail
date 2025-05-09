package orm

import (
	"fmt"
	"net"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
)

func TestAutoMigrate(t *testing.T) {
	t.Run("AutoMigrate", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		_ = AutoMigrate(dbw, &User{}, &Wallet{})
	})
}

func TestPaginate(t *testing.T) {
	t.Run("Paginate", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()

		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		dbw.Scopes(Paginate(1, 100))
	})

	t.Run("Paginate-Condition", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()

		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		dbw.Scopes(Paginate(0, 0))
	})
}

func TestIgnoreErrRecordNotFound(t *testing.T) {
	t.Run("IgnoreErrRecordNotFound", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		_ = IgnoreErrRecordNotFound(dbw)
	})
}
