package nacos

import (
	"testing"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/stretchr/testify/assert"
)

var (
	configJson = []byte(`
{
    "logger_conf": {
        "env": "prod",
        "level": "",
        "filename": "",
        "max_size": 0,
        "max_backups": 0,
        "compress": false,
        "exporter": {
            "provider": "",
            "redis": {
                "list_key": "",
                "conn_conf": {
                    "host": "",
                    "port": 0,
                    "username": "",
                    "password": "",
                    "database": 0,
                    "ssl_enable": false
                },
                "cluster_conn_conf": {
                    "ssl_enable": false,
                    "addr_list": null
                }
            },
            "nats": {
                "subject": "",
                "conn_conf": {
                    "servers": null,
                    "username": "",
                    "password": ""
                }
            }
        }
    }
}`)
	configYaml = []byte(`
logger_conf:
  env: "prod"
  level: ""
  filename: ""
  max_size: 0
  max_backups: 0
  compress: false
  exporter:
    provider: ""
    redis:
      list_key: ""
      conn_conf:
        addr:
          host: ""
          port: 0
          username: ""
          password: ""
        database: 0
        ssl_enable: false
      cluster_conn_conf:
        ssl_enable: false
        addr_list: []
    nats:
      subject: ""
      conn_conf:
        servers: []
        username: ""
        password: ""`)
	configToml = []byte(`
[logger_conf]
env = 'prod'
level = ''
filename = ''
max_size = 0
max_backups = 0
compress = false

[logger_conf.exporter]
provider = ''

[logger_conf.exporter.redis]
list_key = ''

[logger_conf.exporter.redis.conn_conf]
host = ''
port = 0
username = ''
password = ''
database = 0
ssl_enable = false

[logger_conf.exporter.redis.cluster_conn_conf]
ssl_enable = false
addr_list = []

[logger_conf.exporter.nats]
subject = ''

[logger_conf.exporter.nats.conn_conf]
servers = []
username = ''
password = ''`)

	configArr = map[string][]byte{"json": configJson, "toml": configToml, "yaml": configYaml}
)

func TestParseConfig(t *testing.T) {
	for tag, cfg := range configArr {
		var conf config.Config
		err := ParseConfig(cfg, &conf, tag)
		t.Logf("parse format: {%s}, test field env: {%s}, error: %v", tag, conf.LoggerConf.Env, err)
		assert.Equal(t, nil, err)
		assert.Equal(t, "prod", conf.LoggerConf.Env)
	}
}
