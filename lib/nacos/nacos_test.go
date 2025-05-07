package nacos

import (
	"os"
	"testing"

	nacosV2Constant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/stretchr/testify/assert"
)

func TestInitClient(t *testing.T) {
	t.Run("InitClient-NonValue", func(t *testing.T) {
		assert.Panics(t, func() {
			InitClient("go-sail", "", "")
		})
		_ = os.RemoveAll("logs")
	})

	t.Run("InitClient-ServerConfig", func(t *testing.T) {
		conf := nacosV2Constant.ClientConfig{}
		InitClient("go-sail", "127.0.0.1:8848", "abc", conf)
		_ = os.RemoveAll("logs")
	})

	t.Run("InitClient", func(t *testing.T) {
		InitClient("go-sail", "127.0.0.1:8848", "abc")
		_ = os.RemoveAll("logs")
	})
}

func TestGetConfigClient(t *testing.T) {
	t.Run("GetConfigClient", func(t *testing.T) {
		t.Log(GetConfigClient())
	})
}

func TestGetNamingClient(t *testing.T) {
	t.Run("GetNamingClient", func(t *testing.T) {
		t.Log(GetNamingClient())
	})
}
