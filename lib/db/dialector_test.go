package db

import "testing"

func TestMysqlDsn(t *testing.T) {
	t.Run("mysqlDsn", func(t *testing.T) {
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
		t.Log(dbConf.GenDialector())
	})
}
