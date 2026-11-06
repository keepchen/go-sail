package db

import (
	"time"
)

// Conf configuration
//
// <yaml example>
//
// db_conf:
//
//	enable: false
//	driver_name: mysql
//	auto_migrate: true
//	logger:
//	  level: "warn"
//	  slow_threshold: 100
//	  skip_caller_lookup: true
//	  ignore_record_not_found_error: true
//	  colorful: false
//	connection_pool:
//	  # Maximum open connections
//	  max_open_conn_count: 100
//	  # Maximum idle connections
//	  max_idle_conn_count: 10
//	  # Maximum connection lifetime (minutes)
//	  conn_max_life_time_minutes: 30
//	  # Maximum connection idle time (minutes)
//	  conn_max_idle_time_minutes: 10
//	mysql:
//	  read:
//	      host: 127.0.0.1
//	      port: 33060
//	      username: foo
//	      password: bar
//	      database: go_sail
//	      charset: utf8mb4
//	      parseTime: true
//	      loc: Local
//	  write:
//	      host: 127.0.0.1
//	      port: 33060
//	      username: foo
//	      password: bar
//	      database: go_sail
//	      charset: utf8mb4
//	      parseTime: true
//	      loc: Local
//
// <toml example>
//
// # ::Database configuration::
//
// [db_conf]
//
// enable = false
//
// driver_name = "mysql"
//
// # Whether to automatically synchronize table schemas
//
// auto_migrate = false
//
// # Logging configuration
//
// [db_conf.logger]
//
// level = "warn"
//
// slow_threshold = 100
//
// skip_caller_lookup = true
//
// ignore_record_not_found_error = true
//
// colorful = false
//
// # ::Database connection pool configuration::
//
// [db_conf.connection_pool]
//
// # Maximum open connections
//
// max_open_conn_count = 100
//
// # Maximum idle connections
//
// max_idle_conn_count = 10
//
// # Maximum connection lifetime (minutes)
//
// conn_max_life_time_minutes = 30
//
// # Maximum connection idle time (minutes)
//
// conn_max_idle_time_minutes = 10
//
// # MySQL configuration
//
// [db_conf.mysql.read]
//
// host = "localhost"
//
// port = 3306
//
// username = "foo"
//
// password = "bar"
//
// database = "default"
//
// charset = "utf8mb4"
//
// parse_time = true
//
// loc = "Local"
//
// [db_conf.mysql.write]
//
// host = "localhost"
//
// port = 3306
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// charset = "utf8mb4"
//
// parse_time = true
//
// loc = "Local"
//
// # PostgreSQL configuration
//
// [db_conf.postgres.read]
//
// host = "localhost"
//
// port = 9920
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// ssl_mode = "disable" # enable | disable
//
// timezone = "Asia/Shanghai"
//
// # PostgreSQL configuration
//
// [db_conf.postgres.write]
//
// host = "localhost"
//
// port = 9920
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// ssl_mode = "disable" # enable | disable
//
// timezone = "Asia/Shanghai"
//
// # SQL Server configuration
//
// [db_conf.sqlserver.read]
//
// host = "localhost"
//
// port = 9930
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// [db_conf.sqlserver.write]
//
// host = "localhost"
//
// port = 9930
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// # ClickHouse configuration
//
// [db_conf.clickhouse.read]
//
// host = "localhost"
//
// port = 9000
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// read_timeout = 20
//
// write_timeout = 20
//
// [db_conf.clickhouse.write]
//
// host = "localhost"
//
// port = 9000
//
// username = "foo"
//
// password = "bar"
//
// database = "go_sail"
//
// read_timeout = 20
//
// write_timeout = 20
//
// # SQLite configuration
//
// [db_conf.sqlite.read]
//
// file = "sqlite.db"
//
// [db_conf.sqlite.write]
//
// file = "sqlite.db"
type Conf struct {
	Enable                                   bool               `yaml:"enable" toml:"enable" json:"enable"`                                                                                                                      //Whether to enable
	DriverName                               string             `yaml:"driver_name" toml:"driver_name" json:"driver_name"`                                                                                                       //Database type
	AutoMigrate                              bool               `yaml:"auto_migrate" toml:"auto_migrate" json:"auto_migrate"`                                                                                                    //Whether to automatically synchronize table schemas
	DisableForeignKeyConstraintWhenMigrating bool               `yaml:"disable_foreign_key_constraint_when_migrating" toml:"disable_foreign_key_constraint_when_migrating" json:"disable_foreign_key_constraint_when_migrating"` //Disable foreign key constraints during schema migration
	DisableNestedTransaction                 bool               `yaml:"disable_nested_transaction" toml:"disable_nested_transaction" json:"disable_nested_transaction"`                                                          //Disable nested transactions
	AllowGlobalUpdate                        bool               `yaml:"allow_global_update" toml:"allow_global_update" json:"allow_global_update"`                                                                               //Allow global updates (without a WHERE clause)
	SkipDefaultTransaction                   bool               `yaml:"skip_default_transaction" toml:"skip_default_transaction" json:"skip_default_transaction"`                                                                //Skip default transactions
	Logger                                   Logger             `yaml:"logger" toml:"logger" json:"logger"`                                                                                                                      //Logging configuration
	ConnectionPool                           ConnectionPoolConf `yaml:"connection_pool" toml:"connection_pool" json:"connection_pool"`                                                                                           //Connection pool configuration
	Mysql                                    MysqlConf          `yaml:"mysql" toml:"mysql" json:"mysql"`                                                                                                                         //MySQL configuration
	Postgres                                 PostgresConf       `yaml:"postgres" toml:"postgres" json:"postgres"`                                                                                                                //PostgreSQL configuration
	Sqlserver                                SqlserverConf      `yaml:"sqlserver" toml:"sqlserver" json:"sqlserver"`                                                                                                             //SQL Server configuration
	Sqlite                                   SqliteConf         `yaml:"sqlite" toml:"sqlite" json:"sqlite"`                                                                                                                      //SQLite configuration
	Clickhouse                               ClickhouseConf     `yaml:"clickhouse" toml:"clickhouse" json:"clickhouse"`                                                                                                          //ClickHouse configuration
	NowFunc                                  func() time.Time   `yaml:"-" toml:"-" json:"-"`                                                                                                                                     //Current timestamp function @doc: https://gorm.io/docs/gorm_config.html#NowFunc
}

// Logger logging configuration
type Logger struct {
	Level                     string `yaml:"level" toml:"level" json:"level" default:"info"`                                                          //Log level (info,warn,error,silent)
	SlowThreshold             int    `yaml:"slow_threshold" toml:"slow_threshold" json:"slow_threshold" default:"100"`                                //Slow query threshold (milliseconds)
	SkipCallerLookup          bool   `yaml:"skip_caller_lookup" toml:"skip_caller_lookup" json:"skip_caller_lookup"`                                  //Whether to skip caller checks
	IgnoreRecordNotFoundError bool   `yaml:"ignore_record_not_found_error" toml:"ignore_record_not_found_error" json:"ignore_record_not_found_error"` //Ignore record-not-found errors
	Colorful                  bool   `yaml:"colorful" toml:"colorful" json:"colorful"`                                                                //Whether to print in color
}

// ConnectionPoolConf connection pool configuration
type ConnectionPoolConf struct {
	MaxOpenConnCount       int `yaml:"max_open_conn_count" toml:"max_open_conn_count" json:"max_open_conn_count" default:"100"`                     //Maximum open connections
	MaxIdleConnCount       int `yaml:"max_idle_conn_count" toml:"max_idle_conn_count" json:"max_idle_conn_count" default:"10"`                      //Maximum idle connections
	ConnMaxLifeTimeMinutes int `yaml:"conn_max_life_time_minutes" toml:"conn_max_life_time_minutes" json:"conn_max_life_time_minutes" default:"30"` //Maximum connection lifetime (minutes)
	ConnMaxIdleTimeMinutes int `yaml:"conn_max_idle_time_minutes" toml:"conn_max_idle_time_minutes" json:"conn_max_idle_time_minutes" default:"10"` //Maximum connection idle time (minutes)
}

type MysqlConf struct {
	Read  MysqlConfItem `yaml:"read" toml:"read" json:"read"`    //Read instance
	Write MysqlConfItem `yaml:"write" toml:"write" json:"write"` //Write instance
}

// MysqlConfItem MySQL configuration
type MysqlConfItem struct {
	Host      string `yaml:"host" toml:"host" json:"host"`                //Host address
	Port      int    `yaml:"port" toml:"port" json:"port"`                //Port
	Username  string `yaml:"username" toml:"username" json:"username"`    //Username
	Password  string `yaml:"password" toml:"password" json:"password"`    //Password
	Database  string `yaml:"database" toml:"database" json:"database"`    //Database name
	Charset   string `yaml:"charset" toml:"charset" json:"charset"`       //Character set
	ParseTime bool   `yaml:"parseTime" toml:"parseTime" json:"parseTime"` //Whether to parse time values
	Loc       string `yaml:"loc" toml:"loc" json:"loc" default:"Local"`   //Location
}

type PostgresConf struct {
	Read  PostgresConfItem `yaml:"read" toml:"read" json:"read"`    //Read instance
	Write PostgresConfItem `yaml:"write" toml:"write" json:"write"` //Write instance
}

// PostgresConfItem PostgreSQL configuration
type PostgresConfItem struct {
	Host     string `yaml:"host" toml:"host" json:"host"`             //Host address
	Port     int    `yaml:"port" toml:"port" json:"port"`             //Port
	Username string `yaml:"username" toml:"username" json:"username"` //Username
	Password string `yaml:"password" toml:"password" json:"password"` //Password
	Database string `yaml:"database" toml:"database" json:"database"` //Database name
	SSLMode  string `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"` //SSL mode enable|disable
	TimeZone string `yaml:"timezone" toml:"timezone" json:"timezone"` //Time zone
}

type SqlserverConf struct {
	Read  SqlserverConfItem `yaml:"read" toml:"read" json:"read"`    //Read instance
	Write SqlserverConfItem `yaml:"write" toml:"write" json:"write"` //Write instance
}

// SqlserverConfItem SQL Server configuration
type SqlserverConfItem struct {
	Host     string `yaml:"host" toml:"host" json:"host"`             //Host address
	Port     int    `yaml:"port" toml:"port" json:"port"`             //Port
	Username string `yaml:"username" toml:"username" json:"username"` //Username
	Password string `yaml:"password" toml:"password" json:"password"` //Password
	Database string `yaml:"database" toml:"database" json:"database"` //Database name
}

type SqliteConf struct {
	Read  SqliteConfItem `yaml:"read" toml:"read" json:"read"`    //Read instance
	Write SqliteConfItem `yaml:"write" toml:"write" json:"write"` //Write instance
}

// SqliteConfItem SQLite configuration
type SqliteConfItem struct {
	File string `yaml:"file" toml:"file" json:"file"` //Database file
}

type ClickhouseConf struct {
	Read  ClickhouseConfItem `yaml:"read" toml:"read" json:"read"`    //Read instance
	Write ClickhouseConfItem `yaml:"write" toml:"write" json:"write"` //Write instance
}

// ClickhouseConfItem ClickHouse configuration
type ClickhouseConfItem struct {
	Host         string `yaml:"host" toml:"host" json:"host"`                                         //Host address
	Port         int    `yaml:"port" toml:"port" json:"port"`                                         //Port
	Username     string `yaml:"username" toml:"username" json:"username"`                             //Username
	Password     string `yaml:"password" toml:"password" json:"password"`                             //Password
	Database     string `yaml:"database" toml:"database" json:"database"`                             //Database name
	ReadTimeout  int    `yaml:"read_timeout" toml:"read_timeout" json:"read_timeout" default:"20"`    //Read timeout
	WriteTimeout int    `yaml:"write_timeout" toml:"write_timeout" json:"write_timeout" default:"20"` //Write timeout
}
