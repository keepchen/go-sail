package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

//@doc https://open.feishu.cn/cardkit

const (
	//Payload 请求载荷
	Payload = `{
        "timestamp": "%d",
        "sign": "%s",
        "msg_type": "interactive",
        "card": %s
}`
)

// LarkConf lark or 飞书配置
type LarkConf struct {
	WebhookUrl string `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"` //webhook地址(授权参数)
	SignKey    string `yaml:"sign_key" toml:"sign_key" json:"sign_key"`          //签名
}

// LarkResponseEntity 响应实体
type LarkResponseEntity struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// LarkEmit 发射lark通知
func LarkEmit(conf LarkConf, content string) (LarkResponseEntity, error) {
	headers := map[string]string{"Content-Type": "application/json"}
	timestamp := time.Now().Unix()
	sign, _ := genLarkSign(conf.SignKey, timestamp)
	payload := fmt.Sprintf(Payload, timestamp, sign, content)

	var response LarkResponseEntity
	resp, _, err := utils.HttpClient().SendRequest(http.MethodPost, conf.WebhookUrl, []byte(payload), headers, time.Second*10)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(resp, &response)
	return response, err
}

// 生成lark或飞书签名
func genLarkSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)
	return utils.Base64().Encode([]byte(stringToSign)), nil
}
