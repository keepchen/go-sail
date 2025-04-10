package logger

import "testing"

func TestMarshalInterfaceValue(t *testing.T) {
	t.Run("MarshalInterfaceValue", func(t *testing.T) {
		t.Log(MarshalInterfaceValue(map[string]interface{}{
			"name":       "go-sail",
			"opensource": true,
		}))
	})
}
