package nacos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	t.Run("GetConfig", func(t *testing.T) {
		assert.Panics(t, func() {
			conf := struct{}{}
			formats := []string{"yaml", "json", "toml", "unknown"}
			for _, format := range formats {
				t.Log(GetConfig("go-sail", "tester", conf, format))
			}
		})
	})
}

func TestListenConfig(t *testing.T) {
	t.Run("ListenConfig", func(t *testing.T) {
		assert.Panics(t, func() {
			conf := struct{}{}
			formats := []string{"yaml", "json", "toml", "unknown"}
			for _, format := range formats {
				t.Log(ListenConfig("go-sail", "tester", conf, format, false))
				t.Log(ListenConfig("go-sail", "tester", conf, format, true))
			}
		})
	})
}

func TestListenConfigWithCallback(t *testing.T) {
	t.Run("ListenConfigWithCallback", func(t *testing.T) {
		assert.Panics(t, func() {
			fn := func(namespace, group, dataId, data string) {
				fmt.Println(namespace, group, dataId, data)
			}
			t.Log(ListenConfigWithCallback("go-sail", "tester", fn))
		})
	})
}
