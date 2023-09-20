package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/SherClockHolmes/webpush-go"
)

// GenerateVAPIDKeys 生成web push的公私钥
func GenerateVAPIDKeys() (privateKey string, publicKey string, err error) {
	privateKey, publicKey, err = webpush.GenerateVAPIDKeys()
	return
}

// Payload 消息推送载荷
type Payload struct {
	Title string           `json:"title"`
	Body  string           `json:"body"`
	Icon  string           `json:"icon"`
	Data  PayloadDataField `json:"data"`
	Badge string           `json:"badge"`
}

// PayloadDataField 消息推送载荷-data字段
type PayloadDataField struct {
	URL string `json:"url"`
}

// SendNotification 发送通知
func SendNotification(privateKey, publicKey string, sub webpush.Subscription, subscriberEmail string, body Payload) error {
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
