package config

import (
	"encoding/json"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	"github.com/keepchen/go-sail/v3/lib/etcd"
)

func TestPrintConfig(t *testing.T) {
	formats := [...]string{"json", "toml", "yaml", "unknown"}
	for _, format := range formats {
		PrintTemplateConfig(format)
	}
	t.Log("OK")
}

func TestParseConfigFromString(t *testing.T) {
	formats := [...]string{"json", "toml", "yaml", "unknown"}
	conf := &Config{
		HttpServer: HttpServerConf{
			Debug: true,
			Addr:  ":8000",
		},
		EtcdConf: etcd.Conf{
			Timeout: 10,
		},
	}
	js, jsErr := json.Marshal(conf)
	assert.NoError(t, jsErr)
	c1, c1Err := ParseConfigFromBytes(formats[0], js)
	assert.NoError(t, c1Err)
	assert.Equal(t, true, c1.HttpServer.Debug)
	assert.Equal(t, ":8000", c1.HttpServer.Addr)
	assert.Equal(t, 10, c1.EtcdConf.Timeout)

	tml, tmlErr := toml.Marshal(conf)
	assert.NoError(t, tmlErr)
	c2, c2Err := ParseConfigFromBytes(formats[1], tml)
	assert.NoError(t, c2Err)
	assert.Equal(t, true, c2.HttpServer.Debug)
	assert.Equal(t, ":8000", c2.HttpServer.Addr)
	assert.Equal(t, 10, c2.EtcdConf.Timeout)

	ym, ymErr := yaml.Marshal(conf)
	assert.NoError(t, ymErr)
	c3, c3Err := ParseConfigFromBytes(formats[2], ym)
	assert.NoError(t, c3Err)
	assert.Equal(t, true, c3.HttpServer.Debug)
	assert.Equal(t, ":8000", c3.HttpServer.Addr)
	assert.Equal(t, 10, c3.EtcdConf.Timeout)
}
