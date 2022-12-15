package db

import (
	"fmt"
	"strings"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"

	"gorm.io/gorm"
)

//数据库类型标识
const (
	DriverNameMysql      = "mysql"      //数据库类型标识:mysql
	DriverNamePostgres   = "postgres"   //数据库类型标识:postgres sql
	DriverNameSqlite     = "sqlite"     //数据库类型标识:sqlite
	DriverNameSqlserver  = "sqlserver"  //数据库类型标识:sqlserver
	DriverNameClickhouse = "clickhouse" //数据库类型标识:clickhouse

	//more ...
)

//组装mysql的dsn
//@see https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
func mysqlDsn(conf MysqlConfItem) string {
	if len(conf.Charset) == 0 {
		conf.Charset = "utf8mb4"
	}
	if len(conf.Loc) == 0 {
		conf.Loc = "Local"
	}
	var parseTime string
	if conf.ParseTime {
		parseTime = "True"
	} else {
		parseTime = "False"
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		conf.Username, conf.Password, conf.Host, conf.Port,
		conf.Database, conf.Charset, parseTime, conf.Loc)
}

//组装postgres的dsn
//dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
func postgresDsn(conf PostgresConfItem) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		conf.Host, conf.Username, conf.Password, conf.Database,
		conf.Port, conf.SSLMode, conf.TimeZone)
}

//组装sqlite的dsn
//dsn := "sqlite.db"
func sqliteDsn(conf SqliteConfItem) string {
	return conf.File
}

//组装sqlserver的dsn
//dsn := "sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm"
func sqlserverDsn(conf SqlserverConfItem) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		conf.Username, conf.Password, conf.Host,
		conf.Port, conf.Database)
}

//组装clickhouse的dsn
//dsn := "tcp://localhost:9000?database=gorm&username=gorm&password=gorm&read_timeout=10&write_timeout=20"
func clickhouseDsn(conf ClickhouseConfItem) string {
	return fmt.Sprintf("tcp://%s:%d?database=%s&username=%s&password=%s&read_timeout=%d&write_timeout=%d",
		conf.Host, conf.Port, conf.Database,
		conf.Username, conf.Password,
		conf.ReadTimeout, conf.WriteTimeout)
}

//GenDialector 生成数据库连接方言(驱动)
func (conf Conf) GenDialector() (gorm.Dialector, gorm.Dialector) {
	switch strings.ToLower(conf.DriverName) {
	case DriverNameMysql:
		return mysql.Open(mysqlDsn(conf.Mysql.Read)), mysql.Open(mysqlDsn(conf.Mysql.Write))
	case DriverNamePostgres:
		return postgres.Open(postgresDsn(conf.Postgres.Read)), postgres.Open(postgresDsn(conf.Postgres.Write))
	case DriverNameSqlite:
		return sqlite.Open(sqliteDsn(conf.Sqlite.Read)), sqlite.Open(sqliteDsn(conf.Sqlite.Write))
	case DriverNameSqlserver:
		return sqlserver.Open(sqlserverDsn(conf.Sqlserver.Read)), sqlserver.Open(sqlserverDsn(conf.Sqlserver.Write))
	case DriverNameClickhouse:
		return clickhouse.Open(clickhouseDsn(conf.Clickhouse.Read)), clickhouse.Open(clickhouseDsn(conf.Clickhouse.Write))
	default:
		panic("not supported database driver")
	}
}
