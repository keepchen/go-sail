package utils

import (
	"github.com/SherClockHolmes/webpush-go"
	"testing"
)

func TestGenerateVAPIDKeys(t *testing.T) {
	t.Run("GenerateVAPIDKeys", func(t *testing.T) {
		t.Log(WebPush().GenerateVAPIDKeys())
	})
}

func TestSendNotification(t *testing.T) {
	t.Run("SendNotification", func(t *testing.T) {
		pri, pub, err := WebPush().GenerateVAPIDKeys()
		t.Log(pri, pub, err)
		t.Log(WebPush().SendNotification(pri, pub, webpush.Subscription{}, "", Payload{}))
	})
}
