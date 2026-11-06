package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// SlackConf Slack configuration
type SlackConf struct {
	WebhookUrl string `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"` //Webhook URL
	BotToken   string `yaml:"bot_token" toml:"bot_token" json:"bot_token"`       //Bot token; note: ensure the Bot Token has the chat:write permission
	ChannelID  string `yaml:"channel_id" toml:"channel_id" json:"channel_id"`    //Channel
}

// SlackResponseEntity response entity
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

// SlackEmit sends a Slack notification
func SlackEmit(conf SlackConf, content string) (SlackResponseEntity, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", conf.BotToken),
	}
	message := map[string]any{
		"channel": conf.ChannelID,
		"text":    content,
	}
	payload, _ := json.Marshal(message)
	var response SlackResponseEntity
	resp, _, err := utils.HttpClient().SendRequest(http.MethodPost, conf.WebhookUrl, payload, headers, time.Second*10)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	return response, err
}
