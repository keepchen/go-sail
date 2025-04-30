package notification

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestPrintDingTalkConf(t *testing.T) {
	t.Run("PrintDingTalkConf", func(t *testing.T) {
		conf := DingTalkConf{WebhookUrl: "127.0.0.1", Secret: "secret"}

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

func TestParseTalkConf(t *testing.T) {
	t.Run("ParseTalkConf", func(t *testing.T) {
		conf := DingTalkConf{WebhookUrl: "127.0.0.1", Secret: "secret"}

		js, err := json.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(js))

		var jsConf DingTalkConf
		err = json.Unmarshal(js, &jsConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, jsConf.WebhookUrl)
		assert.Equal(t, conf.Secret, jsConf.Secret)

		ym, err := yaml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(ym))

		var ymConf DingTalkConf
		err = yaml.Unmarshal(ym, &ymConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, ymConf.WebhookUrl)
		assert.Equal(t, conf.Secret, ymConf.Secret)

		tm, err := toml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(tm))

		var tmConf DingTalkConf
		err = toml.Unmarshal(tm, &tmConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, tmConf.WebhookUrl)
		assert.Equal(t, conf.Secret, tmConf.Secret)
	})
}

func TestGenDingTalkSign(t *testing.T) {
	t.Run("genDingTalkSign", func(t *testing.T) {
		var (
			secret    = "secret"
			timestamp = time.Now().UnixMilli()
		)
		sign, err := genDingTalkSign(secret, timestamp)
		assert.NoError(t, err)
		t.Log(sign)
	})
}

func TestDingTalkEmit(t *testing.T) {
	t.Run("DingTalkEmit", func(t *testing.T) {
		conf := DingTalkConf{}

		ent, err := DingTalkEmit(conf, "tester-DingTalkEmit")
		t.Log(err)
		t.Log(ent)
	})
}
