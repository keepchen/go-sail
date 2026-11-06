package config

import (
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/jwt"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
)

// Config overall configuration
type Config struct {
	AppName      string            `yaml:"app_name" toml:"app_name" json:"app_name"`                //Application name
	Timezone     string            `yaml:"timezone" toml:"timezone" json:"timezone"`                //Server time zone
	Debug        bool              `yaml:"debug" toml:"debug" json:"debug"`                         //Whether debug mode is enabled
	Logger       logger.Conf       `yaml:"logger" toml:"logger" json:"logger"`                      //Logging
	Datasource   db.Conf           `yaml:"datasource" toml:"datasource" json:"datasource"`          //Database configuration
	Redis        redis.Conf        `yaml:"redis" toml:"redis" json:"redis"`                         //Redis configuration
	RedisCluster redis.ClusterConf `yaml:"redis_cluster" toml:"redis_cluster" json:"redis_cluster"` //Redis cluster configuration
	Nats         nats.Conf         `yaml:"nats" toml:"nats" json:"nats"`                            //NATS configuration
	JWT          jwt.Conf          `yaml:"jwt" toml:"jwt" json:"jwt"`                               //JWT configuration
	HttpServer   HttpServerConf    `yaml:"http_server" toml:"http_server" json:"http_server"`       //HTTP server configuration
}

// HttpServerConf HTTP server configuration
type HttpServerConf struct {
	Addr          string         `yaml:"addr" toml:"addr" json:"addr" default:":8080"`                               //Listen address
	EnableSwagger bool           `yaml:"enable_swagger" toml:"enable_swagger" json:"enable_swagger" default:"false"` //Whether to enable Swagger documentation
	Prometheus    PrometheusConf `yaml:"prometheus_conf" toml:"prometheus_conf" json:"prometheus_conf"`              //Prometheus configuration
}

type PrometheusConf struct {
	Enable bool   `yaml:"enable" toml:"enable" json:"enable" default:"false"` //Whether to enable
	Addr   string `yaml:"addr" toml:"addr" json:"addr" default:":19100"`      //Listen address
}

var globalConfig = &Config{}

// GetGlobalConfig returns the global configuration
func GetGlobalConfig() *Config {
	return globalConfig
}
