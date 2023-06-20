package logger

import (
	"github.com/keepchen/go-sail/v2/pkg/lib/nats"
	"github.com/keepchen/go-sail/v2/pkg/lib/redis"
)

// Conf 日志配置
//
// <yaml example>
//
// logger:
//
//	env: dev
//	level: debug
//	filename: logs/user_running.log
//	max_size: 100
//	max_backups: 10
//	compress: true
//	enable_elk_with_redis_list: false
//	redis_list_key: goal-sail:logs/user_running.log
//
// <toml example>
//
// # ::zap日志组件配置::
//
// [logger]
//
// # 日志环境 dev,prod
//
// env = "dev"
//
// # 日志级别 debug,info,warn,error,dpanic,panic,fatal
//
// level = "info"
//
// # 日志文件名称 需要跟上路径
//
// filename = "logs/running.log"
//
// # 单文件日志大小限制，单位MB
//
// max_size = 100
//
// # 最大历史文件保留数量
//
// max_backups = 10
//
// # 是否压缩历史文件
//
// compress = true
//
// # 是否启用基于redis list的elk日志写入
//
// enable_elk_with_redis_list = true
//
// # redis list的elk日志写入的key
//
// redis_list_key = ""
type Conf struct {
	Env                    string `yaml:"env" toml:"env" json:"env" default:"prod"`                                                                       //日志环境，prod：生产环境，dev：开发环境
	Level                  string `yaml:"level" toml:"level" json:"level" default:"info"`                                                                 //日志级别，debug，info，warn，error
	Filename               string `yaml:"filename" toml:"filename" json:"filename" default:"logs/running.log"`                                            //日志文件名称
	MaxSize                int    `yaml:"max_size" toml:"max_size" json:"max_size" default:"100"`                                                         //日志大小限制，单位MB
	MaxBackups             int    `yaml:"max_backups" toml:"max_backups" json:"max_backups" default:"10"`                                                 //最大历史文件保留数量
	Compress               bool   `yaml:"compress" toml:"compress" json:"compress" default:"true"`                                                        //是否压缩历史日志文件
	EnableELKWithRedisList bool   `yaml:"enable_elk_with_redis_list" toml:"enable_elk_with_redis_list" json:"enable_elk_with_redis_list" default:"false"` //是否启用基于redis list的elk日志写入
	RedisListKey           string `yaml:"redis_list_key" toml:"redis_list_key" json:"redis_list_key"`                                                     //redis list的elk日志写入的key
}

// ConfV2 日志配置
//
// <yaml example>
//
// loggerV2:
//
//	env: dev
//	level: debug
//	filename: logs/user_running.log
//	max_size: 100
//	max_backups: 10
//	compress: true
//	exporter:
//	  provider: "redis-cluster"
//	  nats:
//	    subject: "logger"
//	    conn_conf:
//	        servers:
//	          - "nats://192.168.134.116:4222"
//	        username: admin
//	        password: changeme
//	  redis:
//	    list_key: "go-sail-user:logger"
//	    conn_conf:
//	        addr:
//	          host: ""
//	          port: 0
//	          username: ""
//	          password: ""
//	        database: 0
//	        ssl_enable: false
//	    cluster_conn_conf:
//	      ssl_enable: false
//	      addr_list:
//	        - host: 192.168.224.114
//	          port: 6379
//	          username: ""
//	          password: 123456
//	        - host: 192.168.224.114
//	          port: 6380
//	          username: ""
//	          password: 123456
//
// <toml example>
//
// # ::zap日志组件配置 v2::
//
// [loggerV2]
//
// # 日志环境 dev,prod
//
// env = "dev"
//
// # 日志级别 debug,info,warn,error,dpanic,panic,fatal
//
// level = "info"
//
// # 日志文件名称 需要跟上路径
//
// filename = "logs/running.log"
//
// # 单文件日志大小限制，单位MB
//
// max_size = 100
//
// # 最大历史文件保留数量
//
// max_backups = 10
//
// # 是否压缩历史文件
//
// compress = true
//
// # 日志导出器配置
//
// [loggerV2.exporter]
//
// # 日志导出器介质
//
// provider = "redis"
//
// # nats导出器配置
//
// [loggerV2.exporter.nats]
//
// # nats主题
//
// subject = "logger"
//
// # redis导出器配置
//
// [loggerV2.exporter.redis]
//
// # list键名
//
// list_key = "logger"
//
// [loggerV2.exporter.redis.conn_conf]
//
// host = "localhost"
//
// username = ""
//
// port = 6379
//
// password = ""
//
// database = 0
//
// ssl_enable = false
//
// [loggerV2.exporter.redis.cluster_conn_conf]
//
// [[loggerV2.exporter.redis.cluster_conn_conf.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6380
//
// password = ""
//
// [[loggerV2.exporter.redis.cluster_conn_conf.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6381
//
// password = ""
type ConfV2 struct {
	Env        string `yaml:"env" toml:"env" json:"env" default:"prod"`                            //日志环境，prod：生产环境，dev：开发环境
	Level      string `yaml:"level" toml:"level" json:"level" default:"info"`                      //日志级别，debug，info，warn，error
	Filename   string `yaml:"filename" toml:"filename" json:"filename" default:"logs/running.log"` //日志文件名称
	MaxSize    int    `yaml:"max_size" toml:"max_size" json:"max_size" default:"100"`              //日志大小限制，单位MB
	MaxBackups int    `yaml:"max_backups" toml:"max_backups" json:"max_backups" default:"10"`      //最大历史文件保留数量
	Compress   bool   `yaml:"compress" toml:"compress" json:"compress" default:"true"`             //是否压缩历史日志文件
	Exporter   struct {
		Provider string `yaml:"provider" toml:"provider" json:"provider" default:""` //导出器，目前支持redis、redis-cluster和nats
		Redis    struct {
			ListKey         string            `yaml:"list_key" toml:"list_key" json:"list_key"`                            //redis list的elk日志写入的key
			ConnConf        redis.Conf        `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"`                         //redis连接配置（单机）
			ClusterConnConf redis.ClusterConf `yaml:"cluster_conn_conf" toml:"cluster_conn_conf" json:"cluster_conn_conf"` //redis连接配置（集群）
		} `json:"redis"`
		Nats struct {
			Subject  string    `yaml:"subject" toml:"subject" json:"subject"`       //nats的发布主题
			ConnConf nats.Conf `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"` //nats连接配置
		} `yaml:"nats"`
	} `yaml:"exporter" toml:"exporter" json:"exporter"` //导出器
}
