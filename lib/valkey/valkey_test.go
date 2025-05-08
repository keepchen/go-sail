package valkey

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCaseConf = Conf{
	Endpoints: []Endpoint{
		{Host: "127.0.0.1", Port: 8000},
		{Host: "127.0.0.1", Port: 8001},
		{Host: "127.0.0.1", Port: 8002},
		{Host: "127.0.0.1", Port: 8003},
		{Host: "127.0.0.1", Port: 8004},
		{Host: "127.0.0.1", Port: 8005},
	},
}

func TestInit(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", testCaseConf.Endpoints[0].Host, testCaseConf.Endpoints[0].Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		Init(testCaseConf)
		t.Log("valkey mode:", GetValKey().Mode())
		assert.NoError(t, GetValKey().Do(context.Background(), GetValKey().B().Ping().Build()).Error())

		assert.NoError(t, GetValKey().Do(context.Background(), GetValKey().B().Set().Key("test-Init-set-key").Value("123").Build()).Error())
		resp := GetValKey().Do(context.Background(), GetValKey().B().Get().Key("test-Init-set-key").Build())
		result, err := resp.ToString()
		assert.NoError(t, err)
		assert.Equal(t, "123", result)

		GetValKey().Close()
	})
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", testCaseConf.Endpoints[0].Host, testCaseConf.Endpoints[0].Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		newClient, err := New(testCaseConf)
		assert.NoError(t, err)
		t.Log("valkey mode:", newClient.Mode())
		assert.NoError(t, newClient.Do(context.Background(), newClient.B().Ping().Build()).Error())

		assert.NoError(t, newClient.Do(context.Background(), newClient.B().Set().Key("test-New-set-key").Value("123").Build()).Error())
		resp := newClient.Do(context.Background(), newClient.B().Get().Key("test-New-set-key").Build())
		result, err := resp.ToString()
		assert.NoError(t, err)
		assert.Equal(t, "123", result)

		newClient.Close()
	})
}
