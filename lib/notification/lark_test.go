package notification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestPrintLarkConf(t *testing.T) {
	t.Run("PrintLarkConf", func(t *testing.T) {
		conf := LarkConf{WebhookUrl: "127.0.0.1", SignKey: "secret"}

		js, err := json.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(js))

		ym, err := yaml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(ym))

		tm, err := toml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(tm))
	})
}

func TestParseLarkConf(t *testing.T) {
	t.Run("ParseLarkConf", func(t *testing.T) {
		conf := LarkConf{WebhookUrl: "127.0.0.1", SignKey: "secret"}

		js, err := json.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(js))

		var jsConf LarkConf
		err = json.Unmarshal(js, &jsConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, jsConf.WebhookUrl)
		assert.Equal(t, conf.SignKey, jsConf.SignKey)

		ym, err := yaml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(ym))

		var ymConf LarkConf
		err = yaml.Unmarshal(ym, &ymConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, ymConf.WebhookUrl)
		assert.Equal(t, conf.SignKey, ymConf.SignKey)

		tm, err := toml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(tm))

		var tmConf LarkConf
		err = toml.Unmarshal(tm, &tmConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, tmConf.WebhookUrl)
		assert.Equal(t, conf.SignKey, tmConf.SignKey)
	})
}

func TestGenLarkSign(t *testing.T) {
	t.Run("genLarkSign", func(t *testing.T) {
		var (
			secret    = "secret"
			timestamp = time.Now().Unix()
		)
		sign, err := genLarkSign(secret, timestamp)
		assert.NoError(t, err)
		t.Log(sign)
	})
}

func TestLarkEmit(t *testing.T) {
	t.Run("LarkEmit", func(t *testing.T) {
		conf := LarkConf{}

		ent, err := LarkEmit(conf, "tester-LarkEmit")
		t.Log(err)
		t.Log(ent)
	})
}

func TestLarkEmitPlaintext(t *testing.T) {
	t.Run("LarkEmitPlaintext", func(t *testing.T) {
		conf := LarkConf{}

		ent, err := LarkEmitPlaintext(conf, "tester-LarkEmitPlaintext")
		t.Log(err)
		t.Log(ent)
	})
}
