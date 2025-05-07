package orm

import (
	"fmt"
	"net"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
)

func TestBeforeCreate(t *testing.T) {
	t.Run("BeforeCreate", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		t.Log((&User{}).BeforeCreate(dbw))
	})
}

func TestBeforeSave(t *testing.T) {
	t.Run("BeforeSave", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		t.Log((&User{}).BeforeSave(dbw))
	})
}

func TestBeforeUpdate(t *testing.T) {
	t.Run("BeforeUpdate", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		t.Log((&User{}).BeforeUpdate(dbw))
	})
}

func TestBeforeDelete(t *testing.T) {
	t.Run("BeforeDelete", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dbConf.Mysql.Read.Host, dbConf.Mysql.Read.Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		logger.Init(loggerConf, "go-sail")
		_, _, dbw, _ := db.New(dbConf)
		t.Log((&User{}).BeforeDelete(dbw))
	})
}
