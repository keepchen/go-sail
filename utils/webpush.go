package utils

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/SherClockHolmes/webpush-go"
)

type webPushImpl struct {
}

type IWebPush interface {
	// GenerateVAPIDKeys 生成web push的公私钥
	GenerateVAPIDKeys() (privateKey string, publicKey string, err error)
	// SendNotification 发送通知
	SendNotification(privateKey, publicKey string, sub webpush.Subscription, subscriberEmail string, body Payload) error
}

var wpi IWebPush = &webPushImpl{}

// WebPush 实例化web push工具类
func WebPush() IWebPush {
	return wpi
}

// GenerateVAPIDKeys 生成web push的公私钥
func (webPushImpl) GenerateVAPIDKeys() (privateKey string, publicKey string, err error) {
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
func (webPushImpl) SendNotification(privateKey, publicKey string, sub webpush.Subscription, subscriberEmail string, body Payload) error {
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
