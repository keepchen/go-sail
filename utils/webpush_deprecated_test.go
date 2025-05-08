package utils

import (
	"testing"

	"github.com/SherClockHolmes/webpush-go"
)

func TestGenerateVAPIDKeysDeprecated(t *testing.T) {
	t.Run("GenerateVAPIDKeys", func(t *testing.T) {
		t.Log(GenerateVAPIDKeys())
	})
}

func TestSendNotificationDeprecated(t *testing.T) {
	t.Run("SendNotification", func(t *testing.T) {
		pri, pub, err := GenerateVAPIDKeys()
		t.Log(pri, pub, err)
		t.Log(SendNotification(pri, pub, webpush.Subscription{}, "", PayloadDeprecated{}))
	})
}
