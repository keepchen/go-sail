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

func TestNewConfigClient(t *testing.T) {
	t.Run("NewConfigClient-NonValue", func(t *testing.T) {
		t.Log(NewConfigClient("go-sail", "", ""))
		_ = os.RemoveAll("logs")
	})

	t.Run("NewConfigClient-ServerConfig", func(t *testing.T) {
		conf := nacosV2Constant.ClientConfig{}
		t.Log(NewConfigClient("go-sail", "127.0.0.1:8848", "abc", conf))
		_ = os.RemoveAll("logs")
	})

	t.Run("NewConfigClient", func(t *testing.T) {
		t.Log(NewConfigClient("go-sail", "127.0.0.1:8848", "abc"))
		_ = os.RemoveAll("logs")
	})
}

func TestNewNamingClient(t *testing.T) {
	t.Run("NewNamingClient-NonValue", func(t *testing.T) {
		t.Log(NewNamingClient("go-sail", "", ""))
		_ = os.RemoveAll("logs")
	})

	t.Run("NewNamingClient-ServerConfig", func(t *testing.T) {
		conf := nacosV2Constant.ClientConfig{}
		t.Log(NewNamingClient("go-sail", "127.0.0.1:8848", "abc", conf))
		_ = os.RemoveAll("logs")
	})

	t.Run("NewNamingClient", func(t *testing.T) {
		t.Log(NewNamingClient("go-sail", "127.0.0.1:8848", "abc"))
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
