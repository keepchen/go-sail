package notification

import (
	"encoding/json"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestPrintSlackConf(t *testing.T) {
	t.Run("PrintSlackConf", func(t *testing.T) {
		conf := SlackConf{WebhookUrl: "127.0.0.1", BotToken: "secret"}

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

func TestParseSlackConf(t *testing.T) {
	t.Run("ParseSlackConf", func(t *testing.T) {
		conf := SlackConf{WebhookUrl: "127.0.0.1", BotToken: "secret"}

		js, err := json.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(js))

		var jsConf SlackConf
		err = json.Unmarshal(js, &jsConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, jsConf.WebhookUrl)
		assert.Equal(t, conf.BotToken, jsConf.BotToken)

		ym, err := yaml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(ym))

		var ymConf SlackConf
		err = yaml.Unmarshal(ym, &ymConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, ymConf.WebhookUrl)
		assert.Equal(t, conf.BotToken, ymConf.BotToken)

		tm, err := toml.Marshal(&conf)
		assert.NoError(t, err)
		t.Log(string(tm))

		var tmConf SlackConf
		err = toml.Unmarshal(tm, &tmConf)
		assert.NoError(t, err)
		assert.Equal(t, conf.WebhookUrl, tmConf.WebhookUrl)
		assert.Equal(t, conf.BotToken, tmConf.BotToken)
	})
}

func TestSlackEmit(t *testing.T) {
	t.Run("SlackEmit", func(t *testing.T) {
		conf := SlackConf{}

		ent, err := SlackEmit(conf, "tester-SlackEmit")
		t.Log(err)
		t.Log(ent)
	})
}
