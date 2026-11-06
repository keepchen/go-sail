package config

import (
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/email"
	"github.com/keepchen/go-sail/v3/lib/etcd"
	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/lib/valkey"
)

// Config configuration
type Config struct {
	HttpServer             HttpServerConf    `yaml:"http_conf" toml:"http_conf" json:"http_conf"`                            //HTTP server configuration
	LoggerConf             logger.Conf       `yaml:"logger_conf" toml:"logger_conf" json:"logger_conf"`                      //Logging configuration
	DBConf                 db.Conf           `yaml:"db_conf" toml:"db_conf" json:"db_conf"`                                  //Database configuration
	RedisConf              redis.Conf        `yaml:"redis_conf" toml:"redis_conf" json:"redis_conf"`                         //Redis configuration (standalone)
	RedisClusterConf       redis.ClusterConf `yaml:"redis_cluster_conf" toml:"redis_cluster_conf" json:"redis_cluster_conf"` //Redis configuration (cluster)
	NatsConf               nats.Conf         `yaml:"nats_conf" toml:"nats_conf" json:"nats_conf"`                            //NATS configuration
	JwtConf                *jwt.Conf         `yaml:"jwt_conf" toml:"jwt_conf" json:"jwt_conf"`                               //JWT configuration
	EmailConf              email.Conf        `yaml:"email_conf" toml:"email_conf" json:"email_conf"`                         //Email configuration
	KafkaConf              KafkaExtraConf    `yaml:"kafka_conf" toml:"kafka_conf" json:"kafka_conf"`                         //Kafka configuration
	EtcdConf               etcd.Conf         `yaml:"etcd_conf" toml:"etcd_conf" json:"etcd_conf"`                            //etcd configuration
	ValKeyConf             valkey.Conf       `yaml:"valkey_conf" toml:"valkey_conf" json:"valkey_conf"`                      //Valkey configuration
	MonitorDistributorConf distribution.Conf `yaml:"monitor_conf" toml:"monitor_conf" json:"monitor_conf"`                   //Distributed quantile statistics (pX latency statistics)
}

// HttpServerConf HTTP server configuration
type HttpServerConf struct {
	Debug              bool           `yaml:"debug" toml:"debug" json:"debug" default:"false"`                              //Whether debug mode is enabled
	Addr               string         `yaml:"addr" toml:"addr" json:"addr" default:":8080"`                                 //Listen address
	Swagger            SwaggerConf    `yaml:"swagger_conf" toml:"swagger_conf" json:"swagger_conf"`                         //Swagger documentation configuration
	Prometheus         PrometheusConf `yaml:"prometheus_conf" toml:"prometheus_conf" json:"prometheus_conf"`                //Prometheus configuration
	WebSocketRoutePath string         `yaml:"websocket_route_path" toml:"websocket_route_path" json:"websocket_route_path"` //WebSocket routes
	TrustedProxies     []string       `yaml:"trusted_proxies" toml:"trusted_proxies" json:"trusted_proxies"`                //Trusted proxy addresses, e.g.: 10.0.0.0/16
}

type SwaggerConf struct {
	Enable      bool   `yaml:"enable" toml:"enable" json:"enable" default:"false"`      //Whether to enable
	RedocUIPath string `yaml:"redoc_ui_path" toml:"redoc_ui_path" json:"redoc_ui_path"` //UI page file path, e.g. path/to/docs.html; the filename must be docs.html
	JsonPath    string `yaml:"json_path" toml:"json_path" json:"json_path"`             //JSON file path, e.g. path/to/swagger.json
	FaviconPath string `yaml:"favicon_path" toml:"favicon_path" json:"favicon_path"`    //Browser tab icon file path; defaults to the go-sail icon
}

type PrometheusConf struct {
	Enable              bool   `yaml:"enable" toml:"enable" json:"enable" default:"false"`                                              //Whether to enable
	Addr                string `yaml:"addr" toml:"addr" json:"addr" default:":19100"`                                                   //Listen address
	AccessPath          string `yaml:"access_path" toml:"access_path" json:"access_path" default:"/metrics"`                            //Route path
	DisableSystemSample bool   `yaml:"disable_system_sample" toml:"disable_system_sample" json:"disable_system_sample" default:"false"` //Disable system sampling (enabled by default)
	DiskPath            string `yaml:"disk_path" toml:"disk_path" json:"disk_path" default:"/"`                                         //Monitoring path for checking disk usage
	SampleInterval      string `yaml:"sample_interval" toml:"sample_interval" json:"sample_interval" default:"1m"`                      //Sampling interval (default: 1 minute, minimum: 1ms), format: one minute=1m, 30 seconds=30s, 500 milliseconds=500ms
}

type KafkaExtraConf struct {
	Conf    kafka.Conf `yaml:"conf" toml:"conf" json:"conf"`          //Configuration
	Topic   string     `yaml:"topic" toml:"topic" json:"topic"`       //Subject
	GroupID string     `yaml:"groupID" toml:"groupID" json:"groupID"` //Group ID
}

var config *Config

// Set sets configuration
func Set(conf *Config) {
	config = conf
}

// Get retrieves configuration
func Get() *Config {
	return config
}
