package logger

import (
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
)

// Conf 日志配置
//
// <yaml example>
//
// logger_conf:
//
//	  console_output: false
//	  env: dev
//		 level: debug
//		 modules:
//		   - db
//		   - schedule
//			filename: logs/user_running.log
//			max_size: 100
//			max_backups: 10
//			compress: true
//			exporter:
//			  provider: "redis-cluster"
//			  nats:
//			    subject: "logger"
//			    conn_conf:
//			        servers:
//			          - "nats://192.168.134.116:4222"
//			        username: admin
//			        password: changeme
//			  kafka:
//			    topic: "logger"
//			    conn_conf:
//			        addrList:
//			          - "localhost:9092"
//			        username: admin
//			        password: changeme
//			  redis:
//			    list_key: "go-sail-user:logger"
//			    conn_conf:
//			        endpoint:
//			          host: ""
//			          port: 0
//			          username: ""
//			          password: ""
//			        database: 0
//			        ssl_enable: false
//			    cluster_conn_conf:
//			      ssl_enable: false
//			      endpoints:
//			        - host: 192.168.224.114
//			          port: 6379
//			          username: ""
//			          password: 123456
//			        - host: 192.168.224.114
//			          port: 6380
//			          username: ""
//			          password: 123456
//
// <toml example>
//
// # ::zap日志组件配置 v2::
//
// [logger_conf]
//
// # 是否同时输出到终端
//
// console_output = false
//
// # 日志环境 dev,prod
//
// env = "dev"
//
// # 日志级别 debug,info,warn,error,dpanic,panic,fatal
//
// level = "info"
//
// # 模块名称
//
// modules = ["db", "schedule"]
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
// [logger_conf.exporter]
//
// # 日志导出器介质
//
// provider = "redis"
//
// # nats导出器配置
//
// [logger_conf.exporter.nats]
//
// # nats主题
//
// subject = "logger"
//
// # kafka导出器配置
//
// [logger_conf.exporter.kafka]
//
// # kafka主题
//
// topic = "logger"
//
// # redis导出器配置
//
// [logger_conf.exporter.redis]
//
// # list键名
//
// list_key = "logger"
//
// [logger_conf.exporter.redis.conn_conf]
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
// [logger_conf.exporter.redis.cluster_conn_conf]
//
// [[logger_conf.exporter.redis.cluster_conn_conf.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6380
//
// password = ""
//
// [[logger_conf.exporter.redis.cluster_conn_conf.addr_list]]
//
// host = "localhost"
//
// username = ""
//
// port = 6381
//
// password = ""
type Conf struct {
	ConsoleOutput bool     `yaml:"console_output" toml:"console_output" json:"console_output" default:"false"` //是否同时输出到终端
	Env           string   `yaml:"env" toml:"env" json:"env" default:"prod"`                                   //日志环境，prod：生产环境，dev：开发环境
	Level         string   `yaml:"level" toml:"level" json:"level" default:"info"`                             //日志级别，debug，info，warn，error
	Modules       []string `yaml:"modules" toml:"modules" json:"modules"`                                      //模块名称（日志记录到不同的文件中）
	Filename      string   `yaml:"filename" toml:"filename" json:"filename" default:"logs/running.log"`        //日志文件名称
	MaxSize       int      `yaml:"max_size" toml:"max_size" json:"max_size" default:"100"`                     //日志大小限制，单位MB
	MaxBackups    int      `yaml:"max_backups" toml:"max_backups" json:"max_backups" default:"10"`             //最大历史文件保留数量
	Compress      bool     `yaml:"compress" toml:"compress" json:"compress" default:"true"`                    //是否压缩历史日志文件
	Exporter      Exporter `yaml:"exporter" toml:"exporter" json:"exporter"`                                   //导出器
}

type Exporter struct {
	Provider string        `yaml:"provider" toml:"provider" json:"provider" default:""` //导出器，目前支持redis、redis-cluster、nats和kafka，为空表示不启用
	Redis    ProviderRedis `yaml:"redis" toml:"redis" json:"redis"`
	Nats     ProviderNats  `yaml:"nats" toml:"nats" json:"nats"`
	Kafka    ProviderKafka `yaml:"kafka" toml:"kafka" json:"kafka"`
}

type ProviderRedis struct {
	ListKey         string            `yaml:"list_key" toml:"list_key" json:"list_key"`                            //redis列表名称
	ConnConf        redis.Conf        `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"`                         //redis连接配置（单机）
	ClusterConnConf redis.ClusterConf `yaml:"cluster_conn_conf" toml:"cluster_conn_conf" json:"cluster_conn_conf"` //redis连接配置（集群）
}

type ProviderNats struct {
	Subject  string    `yaml:"subject" toml:"subject" json:"subject"`       //nats的发布主题
	ConnConf nats.Conf `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"` //nats连接配置
}

type ProviderKafka struct {
	Topic    string     `yaml:"topic" toml:"topic" json:"topic"`             //kafka的发布主题
	ConnConf kafka.Conf `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"` //kafka连接配置
}
