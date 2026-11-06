package notification

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

//@doc https://open.feishu.cn/cardkit

const (
	//CardPayload request payload - card
	CardPayload = `{
        "timestamp": "%d",
        "sign": "%s",
        "msg_type": "interactive",
        "card": %s
}`
	//TextPayload request payload - plain text
	TextPayload = `{
        "timestamp": "%d",
        "sign": "%s",
        "msg_type": "text",
        "text": "%s"
}`
)

// LarkConf Lark or Feishu configuration
type LarkConf struct {
	WebhookUrl string `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"` //Webhook URL (authorization parameter)
	SignKey    string `yaml:"sign_key" toml:"sign_key" json:"sign_key"`          //Signature
}

// LarkResponseEntity response entity
type LarkResponseEntity struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// LarkEmit sends a Lark notification - card style
func LarkEmit(conf LarkConf, content string) (LarkResponseEntity, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	timestamp := time.Now().Unix()
	sign, _ := genLarkSign(conf.SignKey, timestamp)
	payload := fmt.Sprintf(CardPayload, timestamp, sign, content)

	var response LarkResponseEntity
	resp, _, err := utils.HttpClient().SendRequest(http.MethodPost, conf.WebhookUrl, []byte(payload), headers, time.Second*10)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	return response, err
}

// LarkEmitPlaintext sends a Lark notification - plain text style
func LarkEmitPlaintext(conf LarkConf, content string) (LarkResponseEntity, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	timestamp := time.Now().Unix()
	sign, _ := genLarkSign(conf.SignKey, timestamp)
	payload := fmt.Sprintf(TextPayload, timestamp, sign, content)

	var response LarkResponseEntity
	resp, _, err := utils.HttpClient().SendRequest(http.MethodPost, conf.WebhookUrl, []byte(payload), headers, time.Second*10)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	return response, err
}

// Generate a Lark or Feishu signature
//
// docs: https://open.larksuite.com/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#3c6592d6
func genLarkSign(secret string, timestamp int64) (string, error) {
	//Apply SHA256 to timestamp + key, then Base64-encode
	stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)

	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}

	return utils.Base64().Encode(h.Sum(nil)), nil
}
