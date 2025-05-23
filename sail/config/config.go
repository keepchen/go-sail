package config

import (
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/email"
	"github.com/keepchen/go-sail/v3/lib/etcd"
	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/lib/kafka"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/lib/valkey"
)

// Config 配置
type Config struct {
	HttpServer       HttpServerConf    `yaml:"http_conf" toml:"http_conf" json:"http_conf"`                            //http服务配置
	LoggerConf       logger.Conf       `yaml:"logger_conf" toml:"logger_conf" json:"logger_conf"`                      //日志配置
	DBConf           db.Conf           `yaml:"db_conf" toml:"db_conf" json:"db_conf"`                                  //数据库配置
	RedisConf        redis.Conf        `yaml:"redis_conf" toml:"redis_conf" json:"redis_conf"`                         //redis配置(standalone)
	RedisClusterConf redis.ClusterConf `yaml:"redis_cluster_conf" toml:"redis_cluster_conf" json:"redis_cluster_conf"` //redis配置(cluster)
	NatsConf         nats.Conf         `yaml:"nats_conf" toml:"nats_conf" json:"nats_conf"`                            //nats配置
	JwtConf          *jwt.Conf         `yaml:"jwt_conf" toml:"jwt_conf" json:"jwt_conf"`                               //jwt配置
	EmailConf        email.Conf        `yaml:"email_conf" toml:"email_conf" json:"email_conf"`                         //邮件配置
	KafkaConf        KafkaExtraConf    `yaml:"kafka_conf" toml:"kafka_conf" json:"kafka_conf"`                         //kafka配置
	EtcdConf         etcd.Conf         `yaml:"etcd_conf" toml:"etcd_conf" json:"etcd_conf"`                            //etcd配置
	ValKeyConf       valkey.Conf       `yaml:"valkey_conf" toml:"valkey_conf" json:"valkey_conf"`                      //valkey配置
}

// HttpServerConf http服务配置
type HttpServerConf struct {
	Debug              bool           `yaml:"debug" toml:"debug" json:"debug" default:"false"`                              //是否是debug模式
	Addr               string         `yaml:"addr" toml:"addr" json:"addr" default:":8080"`                                 //监听地址
	Swagger            SwaggerConf    `yaml:"swagger_conf" toml:"swagger_conf" json:"swagger_conf"`                         //swagger文档配置
	Prometheus         PrometheusConf `yaml:"prometheus_conf" toml:"prometheus_conf" json:"prometheus_conf"`                //prometheus配置
	WebSocketRoutePath string         `yaml:"websocket_route_path" toml:"websocket_route_path" json:"websocket_route_path"` //websocket路由
	TrustedProxies     []string       `yaml:"trusted_proxies" toml:"trusted_proxies" json:"trusted_proxies"`                //受信任的代理地址,如: 10.0.0.0/16
}

type SwaggerConf struct {
	Enable      bool   `yaml:"enable" toml:"enable" json:"enable" default:"false"`      //是否启用
	RedocUIPath string `yaml:"redoc_ui_path" toml:"redoc_ui_path" json:"redoc_ui_path"` //ui页面文件路径，如 path/to/docs.html，注意文件名必须是docs.html
	JsonPath    string `yaml:"json_path" toml:"json_path" json:"json_path"`             //json文件路径，如 path/to/swagger.json
	FaviconPath string `yaml:"favicon_path" toml:"favicon_path" json:"favicon_path"`    //浏览器页签图标文件路径，默认为go-sail图标
}

type PrometheusConf struct {
	Enable              bool   `yaml:"enable" toml:"enable" json:"enable" default:"false"`                                              //是否启用
	Addr                string `yaml:"addr" toml:"addr" json:"addr" default:":19100"`                                                   //监听地址
	AccessPath          string `yaml:"access_path" toml:"access_path" json:"access_path" default:"/metrics"`                            //路由地址
	DisableSystemSample bool   `yaml:"disable_system_sample" toml:"disable_system_sample" json:"disable_system_sample" default:"false"` //禁止系统采样(默认会采样)
	DiskPath            string `yaml:"disk_path" toml:"disk_path" json:"disk_path" default:"/"`                                         //检测硬盘使用率指定的监控路径
	SampleInterval      string `yaml:"sample_interval" toml:"sample_interval" json:"sample_interval" default:"1m"`                      //采样间隔(默认1分钟，最小1ms)，格式：一分钟=1m，30秒=30s，500毫秒=500ms
}

type KafkaExtraConf struct {
	Conf    kafka.Conf `yaml:"conf" toml:"conf" json:"conf"`          //配置
	Topic   string     `yaml:"topic" toml:"topic" json:"topic"`       //主题
	GroupID string     `yaml:"groupID" toml:"groupID" json:"groupID"` //分组id
}

var config *Config

// Set 设置配置
func Set(conf *Config) {
	config = conf
}

// Get 获取配置
func Get() *Config {
	return config
}
