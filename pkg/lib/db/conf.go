package db

//Conf 配置信息
type Conf struct {
	DriverName     string             `yaml:"diver_name" toml:"diver_name" json:"diver_name" default:"mysql"`       //数据库类型
	AutoMigrate    bool               `yaml:"auto_migrate" toml:"auto_migrate" json:"auto_migrate" default:"false"` //是否自动同步表结构
	LogLevel       string             `yaml:"log_level" toml:"log_level" json:"log_level" default:"info"`           //日志级别
	ConnectionPool ConnectionPoolConf `yaml:"connection_pool" toml:"connection_pool" json:"connection_pool"`        //连接池配置
	Mysql          MysqlConf          `yaml:"mysql" toml:"mysql" json:"mysql"`                                      //mysql配置
	Postgres       PostgresConf       `yaml:"postgres" toml:"postgres" json:"postgres"`                             //postgres配置
	Sqlserver      SqlserverConf      `yaml:"sqlserver" toml:"sqlserver" json:"sqlserver"`                          //sqlserver配置
	Sqlite         SqliteConf         `yaml:"sqlite" toml:"sqlite" json:"sqlite"`                                   //sqlite配置
	Clickhouse     ClickhouseConf     `yaml:"clickhouse" toml:"clickhouse" json:"clickhouse"`                       //clickhouse配置
}

//ConnectionPoolConf 连接池配置
type ConnectionPoolConf struct {
	MaxOpenConnCount       int `yaml:"max_open_conn_count" toml:"max_open_conn_count" json:"max_open_conn_count" default:"100"`                     //最大开启连接数
	MaxIdleConnCount       int `yaml:"max_idle_conn_count" toml:"max_idle_conn_count" json:"max_idle_conn_count" default:"10"`                      //最大闲置数量
	ConnMaxLifeTimeMinutes int `yaml:"conn_max_life_time_minutes" toml:"conn_max_life_time_minutes" json:"conn_max_life_time_minutes" default:"30"` //连接最大存活时间(分钟)
	ConnMaxIdleTimeMinutes int `yaml:"conn_max_idle_time_minutes" toml:"conn_max_idle_time_minutes" json:"conn_max_idle_time_minutes" default:"10"` //连接最大空闲时间(分钟)
}

type MysqlConf struct {
	Read  MysqlConfItem `yaml:"read" toml:"read" json:"read" default:"localhost"`    //读实例
	Write MysqlConfItem `yaml:"write" toml:"write" json:"write" default:"localhost"` //写实例
}

//MysqlConfItem mysql配置
type MysqlConfItem struct {
	Host      string `yaml:"host" toml:"host" json:"host" default:"localhost"`           //主机地址
	Port      int    `yaml:"port" toml:"port" json:"port" default:"3306"`                //端口
	Username  string `yaml:"username" toml:"username" json:"username"`                   //用户名
	Password  string `yaml:"password" toml:"password" json:"password"`                   //密码
	Database  string `yaml:"database" toml:"database" json:"database"`                   //数据库名
	Charset   string `yaml:"charset" toml:"charset" json:"charset"`                      //字符集
	ParseTime bool   `yaml:"parseTime" toml:"parseTime" json:"parseTime" default:"true"` //是否解析时间
	Loc       string `yaml:"loc" toml:"loc" json:"loc" default:"Local"`                  //位置
}

type PostgresConf struct {
	Read  PostgresConfItem `yaml:"read" toml:"read" json:"read" default:"localhost"`    //读实例
	Write PostgresConfItem `yaml:"write" toml:"write" json:"write" default:"localhost"` //写实例
}

//PostgresConfItem postgres配置
type PostgresConfItem struct {
	Host     string `yaml:"host" toml:"host" json:"host" default:"localhost"`                 //主机地址
	Port     int    `yaml:"port" toml:"port" json:"port" default:"9920"`                      //端口
	Username string `yaml:"username" toml:"username" json:"username"`                         //用户名
	Password string `yaml:"password" toml:"password" json:"password"`                         //密码
	Database string `yaml:"database" toml:"database" json:"database"`                         //数据库名
	SSLMode  string `yaml:"ssl_mode" toml:"ssl_mode" json:"ssl_mode"`                         //ssl模式 enable|disable
	TimeZone string `yaml:"timezone" toml:"timezone" json:"timezone" default:"Asia/Shanghai"` //时区
}

type SqlserverConf struct {
	Read  SqlserverConfItem `yaml:"read" toml:"read" json:"read" default:"localhost"`    //读实例
	Write SqlserverConfItem `yaml:"write" toml:"write" json:"write" default:"localhost"` //写实例
}

//SqlserverConfItem sqlserver配置
type SqlserverConfItem struct {
	Host     string `yaml:"host" toml:"host" json:"host" default:"localhost"` //主机地址
	Port     int    `yaml:"port" toml:"port" json:"port" default:"9930"`      //端口
	Username string `yaml:"username" toml:"username" json:"username"`         //用户名
	Password string `yaml:"password" toml:"password" json:"password"`         //密码
	Database string `yaml:"database" toml:"database" json:"database"`         //数据库名
}

type SqliteConf struct {
	Read  SqliteConfItem `yaml:"read" toml:"read" json:"read" default:"localhost"`    //读实例
	Write SqliteConfItem `yaml:"write" toml:"write" json:"write" default:"localhost"` //写实例
}

//SqliteConfItem sqlite配置
type SqliteConfItem struct {
	File string `yaml:"file" toml:"file" json:"file" default:"sqlite.db"` //数据库文件
}

type ClickhouseConf struct {
	Read  ClickhouseConfItem `yaml:"read" toml:"read" json:"read" default:"localhost"`    //读实例
	Write ClickhouseConfItem `yaml:"write" toml:"write" json:"write" default:"localhost"` //写实例
}

//ClickhouseConfItem clickhouse配置
type ClickhouseConfItem struct {
	Host         string `yaml:"host" toml:"host" json:"host" default:"localhost"`                     //主机地址
	Port         int    `yaml:"port" toml:"port" json:"port" default:"9000"`                          //端口
	Username     string `yaml:"username" toml:"username" json:"username"`                             //用户名
	Password     string `yaml:"password" toml:"password" json:"password"`                             //密码
	Database     string `yaml:"database" toml:"database" json:"database"`                             //数据库名
	ReadTimeout  int    `yaml:"read_timeout" toml:"read_timeout" json:"read_timeout" default:"20"`    //读取超时时间
	WriteTimeout int    `yaml:"write_timeout" toml:"write_timeout" json:"write_timeout" default:"20"` //写入超时时间
}
