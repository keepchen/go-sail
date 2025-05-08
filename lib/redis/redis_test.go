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
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
	if err != nil {
		return
	}
	_ = conn.Close()

	t.Run("InitRedis", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		InitRedis(sConf)

		assert.Equal(t, true, GetInstance() != nil)

		_, err = GetInstance().Ping(ctx).Result()
		assert.NoError(t, err)

		_, err = GetInstance().Set(ctx, "tester-InitRedis-set", "go-sail", time.Minute).Result()
		assert.NoError(t, err)

		result, err := GetInstance().Get(ctx, "tester-InitRedis-set").Result()
		assert.NoError(t, err)
		t.Log(result)
		assert.Equal(t, "go-sail", result)

		_ = GetInstance().Close()
	})

	t.Run("InitRedis-Panic", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			//unknown port
			sConf.Port += 1
			InitRedis(sConf)

			assert.Equal(t, true, GetInstance() != nil)

			_, err = GetInstance().Ping(ctx).Result()
			assert.NoError(t, err)

			_, err = GetInstance().Set(ctx, "tester-InitRedis-set", "go-sail", time.Minute).Result()
			assert.NoError(t, err)

			result, err := GetInstance().Get(ctx, "tester-InitRedis-set").Result()
			assert.NoError(t, err)
			t.Log(result)
			assert.Equal(t, "go-sail", result)

			_ = GetInstance().Close()
		})
	})

	t.Run("InitRedis-Auth", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			//unknown port
			sConf.Port += 1
			sConf.Username = "username"
			sConf.Password = "password"
			InitRedis(sConf)

			assert.Equal(t, true, GetInstance() != nil)

			_, err = GetInstance().Ping(ctx).Result()
			assert.NoError(t, err)

			_, err = GetInstance().Set(ctx, "tester-InitRedis-set", "go-sail", time.Minute).Result()
			assert.NoError(t, err)

			result, err := GetInstance().Get(ctx, "tester-InitRedis-set").Result()
			assert.NoError(t, err)
			t.Log(result)
			assert.Equal(t, "go-sail", result)

			_ = GetInstance().Close()
		})
	})

	t.Run("InitRedis-SSLEnable", func(t *testing.T) {
		assert.Panics(t, func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			sConf.SSLEnable = true
			InitRedis(sConf)

			assert.Equal(t, true, GetInstance() != nil)

			_, err = GetInstance().Ping(ctx).Result()
			assert.NoError(t, err)

			_, err = GetInstance().Set(ctx, "tester-InitRedis-set", "go-sail", time.Minute).Result()
			assert.NoError(t, err)

			result, err := GetInstance().Get(ctx, "tester-InitRedis-set").Result()
			assert.NoError(t, err)
			t.Log(result)
			assert.Equal(t, "go-sail", result)

			_ = GetInstance().Close()
		})
	})
}

func TestNewRedis(t *testing.T) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
	if err != nil {
		return
	}
	_ = conn.Close()

	t.Run("NewRedis", func(t *testing.T) {
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

	t.Run("NewRedis-SSLEnable", func(t *testing.T) {
		sConf.SSLEnable = true
		newClient, err := New(sConf)
		assert.Error(t, err)
		t.Log(newClient)

		_, err = newClient.Ping(context.Background()).Result()
		assert.Error(t, err)

		_, err = newClient.Set(context.Background(), "tester-NewRedis-set", "go-sail", time.Minute).Result()
		assert.Error(t, err)

		result, err := newClient.Get(context.Background(), "tester-NewRedis-set").Result()
		assert.Error(t, err)
		t.Log(result)

		_ = newClient.Close()
	})

	t.Run("NewRedis-Auth", func(t *testing.T) {
		sConf.SSLEnable = true
		sConf.Username = "username"
		sConf.Password = "password"
		newClient, err := New(sConf)
		assert.Error(t, err)
		t.Log(newClient)

		_, err = newClient.Ping(context.Background()).Result()
		assert.Error(t, err)

		_, err = newClient.Set(context.Background(), "tester-NewRedis-set", "go-sail", time.Minute).Result()
		assert.Error(t, err)

		result, err := newClient.Get(context.Background(), "tester-NewRedis-set").Result()
		assert.Error(t, err)
		t.Log(result)

		_ = newClient.Close()
	})
}
