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
	// GenerateVAPIDKeys generates public and private keys for Web Push
	GenerateVAPIDKeys() (privateKey string, publicKey string, err error)
	// SendNotification sends a notification
	SendNotification(privateKey, publicKey string, sub webpush.Subscription, subscriberEmail string, body Payload) error
}

var wpi IWebPush = &webPushImpl{}

// WebPush instantiates Web Push utilities
func WebPush() IWebPush {
	return wpi
}

// GenerateVAPIDKeys generates public and private keys for Web Push
func (webPushImpl) GenerateVAPIDKeys() (privateKey string, publicKey string, err error) {
	privateKey, publicKey, err = webpush.GenerateVAPIDKeys()
	return
}

// Payload push notification payload
type Payload struct {
	Title string           `json:"title"`
	Body  string           `json:"body"`
	Icon  string           `json:"icon"`
	Data  PayloadDataField `json:"data"`
	Badge string           `json:"badge"`
}

// PayloadDataField push notification payload - data field
type PayloadDataField struct {
	URL string `json:"url"`
}

// SendNotification sends a notification
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
	defer func() {
		_ = respData.Body.Close()
	}()
	if respData.StatusCode < 200 || respData.StatusCode > 299 {
		return fmt.Errorf("endpoint: %s | response body: %s | error: %v", sub.Endpoint, string(respBody), respErr)
	}

	return nil
}
