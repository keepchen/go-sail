package templates

import (
	"strings"
)

// AppConfigTpl 应用配置模板
var AppConfigTpl = strings.ReplaceAll(`package config

import (
	"github.com/keepchen/go-sail/v3/example/pkg/lib/db"
	"github.com/keepchen/go-sail/v3/example/pkg/lib/jwt"
	"github.com/keepchen/go-sail/v3/example/pkg/lib/logger"
	"github.com/keepchen/go-sail/v3/example/pkg/lib/redis"
)

// Config 整体的配置信息
type Config struct {
	AppName      string            <->yaml:"app_name" toml:"app_name" json:"app_name"<->                //应用名称
	Timezone     string            <->yaml:"timezone" toml:"timezone" json:"timezone"<->                //服务器时区
	Debug        bool              <->yaml:"debug" toml:"debug" json:"debug"<->                         //是否是调试模式
	Logger       logger.Conf       <->yaml:"logger" toml:"logger" json:"logger"<->                      //日志
	Datasource   db.Conf           <->yaml:"datasource" toml:"datasource" json:"datasource"<->          //数据库配置
	Redis        redis.Conf        <->yaml:"redis" toml:"redis" json:"redis"<->                         //redis配置
	RedisCluster redis.ClusterConf <->yaml:"redis_cluster" toml:"redis_cluster" json:"redis_cluster"<-> //redis集群配置
	JWT          jwt.Conf          <->yaml:"jwt" toml:"jwt" json:"jwt"<->                               //jwt配置
	HttpServer   HttpServerConf    <->yaml:"http_server" toml:"http_server" json:"http_server"<->       //http服务配置
}

// HttpServerConf http服务配置
type HttpServerConf struct {
	Addr          string         <->yaml:"addr" toml:"addr" json:"addr" default:":8080"<->                               //监听地址
	EnableSwagger bool           <->yaml:"enable_swagger" toml:"enable_swagger" json:"enable_swagger" default:"false"<-> //是否开启swagger文档
	Prometheus    PrometheusConf <->yaml:"prometheus_conf" toml:"prometheus_conf" json:"prometheus_conf"<->              //prometheus配置
}

type PrometheusConf struct {
	Enable bool   <->yaml:"enable" toml:"enable" json:"enable" default:"false"<-> //是否启用
	Addr   string <->yaml:"addr" toml:"addr" json:"addr" default:":19100"<->      //监听地址
}

var globalConfig = &Config{}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	return globalConfig
}
`, "<->", "`")
