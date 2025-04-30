package valkey

import (
	"encoding/json"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintConf(t *testing.T) {
	t.Run("PrintConf", func(t *testing.T) {
		conf := Conf{}
		js, err := json.Marshal(&conf)
		t.Log(string(js))
		assert.NoError(t, err)

		tm, err := toml.Marshal(&conf)
		t.Log(string(tm))
		assert.NoError(t, err)

		ym, err := yaml.Marshal(&conf)
		t.Log(string(ym))
		assert.NoError(t, err)
	})
}

func TestParseConf(t *testing.T) {
	conf := Conf{
		Enable:   true,
		Username: "username",
		Password: "password",
	}
	js, err := json.Marshal(&conf)
	t.Log(string(js))
	assert.NoError(t, err)

	var jsConf Conf
	err = json.Unmarshal(js, &jsConf)
	assert.NoError(t, err)
	assert.Equal(t, conf.Enable, jsConf.Enable)
	assert.Equal(t, conf.Username, jsConf.Username)
	assert.Equal(t, conf.Password, jsConf.Password)

	tm, err := toml.Marshal(&conf)
	t.Log(string(tm))
	assert.NoError(t, err)

	var tmConf Conf
	err = toml.Unmarshal(tm, &tmConf)
	assert.NoError(t, err)
	assert.Equal(t, conf.Enable, tmConf.Enable)
	assert.Equal(t, conf.Username, tmConf.Username)
	assert.Equal(t, conf.Password, tmConf.Password)

	ym, err := yaml.Marshal(&conf)
	t.Log(string(ym))
	assert.NoError(t, err)

	var ymConf Conf
	err = yaml.Unmarshal(ym, &ymConf)
	assert.NoError(t, err)
	assert.Equal(t, conf.Enable, ymConf.Enable)
	assert.Equal(t, conf.Username, ymConf.Username)
	assert.Equal(t, conf.Password, ymConf.Password)
}
