package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// DingTalkConf 钉钉或DingTalk配置
type DingTalkConf struct {
	WebhookUrl string `yaml:"webhook_url" toml:"webhook_url" json:"webhook_url"` //webhook地址(授权参数)
	Secret     string `yaml:"secret" toml:"secret" json:"secret"`                //密钥
}

// DingTalkResponseEntity 响应实体
type DingTalkResponseEntity struct {
	ErrorCode    int    `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

// DingTalkEmit 发射DingTalk通知
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

// 生成DingTalk或钉钉签名
func genDingTalkSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)
	return utils.Base64().Encode([]byte(stringToSign)), nil
}
