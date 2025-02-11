package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/SherClockHolmes/webpush-go"
)

// GenerateVAPIDKeys 生成web push的公私钥
//
// Deprecated: GenerateVAPIDKeys is deprecated,it will be removed in the future.
//
// Please use WebPush().GenerateVAPIDKeys() instead.
func GenerateVAPIDKeys() (privateKey string, publicKey string, err error) {
	privateKey, publicKey, err = webpush.GenerateVAPIDKeys()
	return
}

// PayloadDeprecated 消息推送载荷
type PayloadDeprecated struct {
	Title string                     `json:"title"`
	Body  string                     `json:"body"`
	Icon  string                     `json:"icon"`
	Data  PayloadDataFieldDeprecated `json:"data"`
	Badge string                     `json:"badge"`
}

// PayloadDataFieldDeprecated 消息推送载荷-data字段
type PayloadDataFieldDeprecated struct {
	URL string `json:"url"`
}

// SendNotification 发送通知
//
// Deprecated: SendNotification is deprecated,it will be removed in the future.
//
// Please use WebPush().SendNotification() instead.
func SendNotification(privateKey, publicKey string, sub webpush.Subscription, subscriberEmail string, body PayloadDeprecated) error {
	js, _ := json.Marshal(body)
	respData, respErr := webpush.SendNotification(js, &sub, &webpush.Options{
		Subscriber:      subscriberEmail,
		VAPIDPrivateKey: privateKey,
		VAPIDPublicKey:  publicKey,
		TTL:             30,
	})
	if respData == nil {
		return fmt.Errorf("endpoint: %s | response body: %s | error: %v", sub.Endpoint, "", respErr)
	}
	respBody, _ := io.ReadAll(respData.Body)
	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return fmt.Errorf("endpoint: %s | response body: %s | error: %v", sub.Endpoint, string(respBody), respErr)
	}

	return nil
}
