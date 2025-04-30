package redis

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var sConf = Conf{
	Enable: true,
	Endpoint: Endpoint{
		Host:     "127.0.0.1",
		Port:     6379,
		Username: "",
		Password: "",
	},
	Database: 0,
}

func TestInitRedis(t *testing.T) {
	t.Run("InitRedis", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
		if err != nil {
			return
		}
		_ = conn.Close()

		InitRedis(sConf)

		assert.Equal(t, true, GetInstance() != nil)

		_, err = GetInstance().Ping(context.Background()).Result()
		assert.NoError(t, err)

		_, err = GetInstance().Set(context.Background(), "tester-InitRedis-set", "go-sail", time.Minute).Result()
		assert.NoError(t, err)

		result, err := GetInstance().Get(context.Background(), "tester-InitRedis-set").Result()
		assert.NoError(t, err)
		t.Log(result)
		assert.Equal(t, "go-sail", result)

		_ = GetInstance().Close()
	})
}

func TestNewRedis(t *testing.T) {
	t.Run("NewRedis", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
		if err != nil {
			return
		}
		_ = conn.Close()

		newClient, err := New(sConf)
		assert.NoError(t, err)
		assert.Equal(t, true, newClient != nil)

		_, err = newClient.Ping(context.Background()).Result()
		assert.NoError(t, err)

		_, err = newClient.Set(context.Background(), "tester-NewRedis-set", "go-sail", time.Minute).Result()
		assert.NoError(t, err)

		result, err := newClient.Get(context.Background(), "tester-NewRedis-set").Result()
		assert.NoError(t, err)
		t.Log(result)
		assert.Equal(t, "go-sail", result)

		_ = newClient.Close()
	})
}
