package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// DingTalkConf DingTalk configuration
type DingTalkConf struct {
	WebhookUrl string `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"` //Webhook URL (authorization parameter)
	Secret     string `yaml:"secret" toml:"secret" json:"secret"`                //Secret key
}

// DingTalkResponseEntity response entity
type DingTalkResponseEntity struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

// DingTalkEmit sends a DingTalk notification
func DingTalkEmit(conf DingTalkConf, content string) (DingTalkResponseEntity, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	timestamp := time.Now().UnixMilli()
	sign, _ := genDingTalkSign(conf.Secret, timestamp)

	var response DingTalkResponseEntity
	url := fmt.Sprintf("%s?timestamp=%d&sign=%s", conf.WebhookUrl, timestamp, sign)
	resp, _, err := utils.HttpClient().SendRequest(http.MethodPost, url, []byte(content), headers, time.Second*10)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	return response, err
}

// Generate a DingTalk signature
func genDingTalkSign(secret string, timestamp int64) (string, error) {
	//Apply SHA256 to timestamp + key, then Base64-encode
	stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)
	return utils.Base64().Encode([]byte(stringToSign)), nil
}
