package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// SlackConf slack配置
type SlackConf struct {
	WebhookUrl string `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"` //webhook地址
	BotToken   string `yaml:"bot_token" toml:"bot_token" json:"bot_token"`       //bot token，注意：确保Bot Token 拥有 chat:write 权限
	ChannelID  string `yaml:"channel_id" toml:"channel_id" json:"channel_id"`    //频道
}

// SlackResponseEntity 响应实体
type SlackResponseEntity struct {
	OK        bool   `json:"ok"`
	Channel   string `json:"channel"`
	Timestamp string `json:"ts"` //eg. "1640968740.000200"
	Message   struct {
		Text      string `json:"text"`
		User      string `json:"user"`
		Timestamp string `json:"ts"` //eg. "1640968740.000200"
	} `json:"message"`
}

// SlackEmit 发射slack通知
func SlackEmit(conf SlackConf, content string) (SlackResponseEntity, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", conf.BotToken),
	}
	message := map[string]interface{}{
		"channel": conf.ChannelID,
		"text":    content,
	}
	payload, _ := json.Marshal(message)
	var response SlackResponseEntity
	resp, _, err := utils.SendRequest(http.MethodPost, conf.WebhookUrl, payload, headers, time.Second*10)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	return response, err
}
