package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInterfaceValue(t *testing.T) {
	t.Run("MarshalInterfaceValue", func(t *testing.T) {
		assert.NotNil(t, MarshalInterfaceValue(map[string]any{
			"name":       "go-sail",
			"opensource": true,
		}))
	})
}
