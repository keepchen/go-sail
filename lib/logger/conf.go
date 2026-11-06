package logger

import (
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
)

// Conf logging configuration
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
// # ::Zap logging component configuration v2::
//
// [logger_conf]
//
// # Whether to also output to the terminal
//
// console_output = false
//
// # Logging environment dev,prod
//
// env = "dev"
//
// # Log level debug,info,warn,error,dpanic,panic,fatal
//
// level = "info"
//
// # Module names
//
// modules = ["db", "schedule"]
//
// # Log filename, including the path
//
// filename = "logs/running.log"
//
// # Maximum size per log file, in MB
//
// max_size = 100
//
// # Maximum number of old log files to retain
//
// max_backups = 10
//
// # Whether to compress old log files
//
// compress = true
//
// # Log exporter configuration
//
// [logger_conf.exporter]
//
// # Log exporter transport
//
// provider = "redis"
//
// # NATS exporter configuration
//
// [logger_conf.exporter.nats]
//
// # NATS subject
//
// subject = "logger"
//
// # Kafka exporter configuration
//
// [logger_conf.exporter.kafka]
//
// # Kafka topic
//
// topic = "logger"
//
// # Redis exporter configuration
//
// [logger_conf.exporter.redis]
//
// # List key name
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
	ConsoleOutput bool     `yaml:"console_output" toml:"console_output" json:"console_output" default:"false"` //Whether to also output to the terminal
	Env           string   `yaml:"env" toml:"env" json:"env" default:"prod"`                                   //Logging environment, prod: production, dev: development
	Level         string   `yaml:"level" toml:"level" json:"level" default:"debug"`                            //Log level, debug, info, warn, error
	Modules       []string `yaml:"modules" toml:"modules" json:"modules"`                                      //Module names (logs are written to separate files)
	Filename      string   `yaml:"filename" toml:"filename" json:"filename" default:"logs/running.log"`        //Log filename
	MaxSize       int      `yaml:"max_size" toml:"max_size" json:"max_size" default:"100"`                     //Log size limit, in MB
	MaxBackups    int      `yaml:"max_backups" toml:"max_backups" json:"max_backups" default:"10"`             //Maximum number of old log files to retain
	Compress      bool     `yaml:"compress" toml:"compress" json:"compress"`                                   //Whether to compress old log files
	Exporter      Exporter `yaml:"exporter" toml:"exporter" json:"exporter"`                                   //Exporter
}

type Exporter struct {
	Provider string        `yaml:"provider" toml:"provider" json:"provider"` //Exporter; supports redis, redis-cluster, nats, and kafka; empty means disabled
	Redis    ProviderRedis `yaml:"redis" toml:"redis" json:"redis"`
	Nats     ProviderNats  `yaml:"nats" toml:"nats" json:"nats"`
	Kafka    ProviderKafka `yaml:"kafka" toml:"kafka" json:"kafka"`
}

type ProviderRedis struct {
	ListKey         string            `yaml:"list_key" toml:"list_key" json:"list_key"`                            //Redis list name
	ConnConf        redis.Conf        `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"`                         //Redis connection configuration (standalone)
	ClusterConnConf redis.ClusterConf `yaml:"cluster_conn_conf" toml:"cluster_conn_conf" json:"cluster_conn_conf"` //Redis connection configuration (cluster)
}

type ProviderNats struct {
	Subject  string    `yaml:"subject" toml:"subject" json:"subject"`       //NATS publish subject
	ConnConf nats.Conf `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"` //NATS connection configuration
}

type ProviderKafka struct {
	Topic    string     `yaml:"topic" toml:"topic" json:"topic"`             //Kafka publish topic
	ConnConf kafka.Conf `yaml:"conn_conf" toml:"conn_conf" json:"conn_conf"` //Kafka connection configuration
}
