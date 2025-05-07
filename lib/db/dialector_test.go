package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMysqlDsn(t *testing.T) {
	t.Run("mysqlDsn", func(t *testing.T) {
		t.Log(mysqlDsn(dbConf.Mysql.Read))
	})

	t.Run("mysqlDsn-NoneValue", func(t *testing.T) {
		dbConf.Mysql.Read.Charset = ""
		dbConf.Mysql.Read.Loc = ""
		dbConf.Mysql.Read.ParseTime = false
		t.Log(mysqlDsn(dbConf.Mysql.Read))
	})
}

func TestPostgresDsn(t *testing.T) {
	t.Run("postgresDsn", func(t *testing.T) {
		t.Log(postgresDsn(dbConf.Postgres.Read))
	})
}

func TestSqliteDsn(t *testing.T) {
	t.Run("sqliteDsn", func(t *testing.T) {
		t.Log(sqliteDsn(dbConf.Sqlite.Read))
	})
}

func TestSqlserverDsn(t *testing.T) {
	t.Run("sqlserverDsn", func(t *testing.T) {
		t.Log(sqlserverDsn(dbConf.Sqlserver.Read))
	})
}

func TestClickhouseDsn(t *testing.T) {
	t.Run("clickhouseDsn", func(t *testing.T) {
		t.Log(clickhouseDsn(dbConf.Clickhouse.Read))
	})
}

func TestGenDialector(t *testing.T) {
	t.Run("GenDialector", func(t *testing.T) {
		dbConf.DriverName = DriverNameMysql
		t.Log(dbConf.GenDialector())
		dbConf.DriverName = DriverNamePostgres
		t.Log(dbConf.GenDialector())
		dbConf.DriverName = DriverNameSqlite
		t.Log(dbConf.GenDialector())
		dbConf.DriverName = DriverNameSqlserver
		t.Log(dbConf.GenDialector())
		dbConf.DriverName = DriverNameClickhouse
		t.Log(dbConf.GenDialector())
	})

	t.Run("GenDialector-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			dbConf.DriverName = "UnknownDriver"
			dbConf.GenDialector()
		})
	})
}
