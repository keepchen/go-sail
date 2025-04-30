package redis

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var cConf = ClusterConf{
	Enable: true,
	Endpoints: []Endpoint{
		{Host: "192.168.30.2", Port: 6379, Username: "", Password: "Tqzk12356"},
		{Host: "192.168.30.2", Port: 6380, Username: "", Password: "Tqzk12356"},
		{Host: "192.168.30.2", Port: 6381, Username: "", Password: "Tqzk12356"},
		{Host: "192.168.30.2", Port: 6382, Username: "", Password: "Tqzk12356"},
		{Host: "192.168.30.2", Port: 6383, Username: "", Password: "Tqzk12356"},
		{Host: "192.168.30.2", Port: 6384, Username: "", Password: "Tqzk12356"},
	},
}

func TestInitRedisCluster(t *testing.T) {
	t.Run("InitRedisCluster", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cConf.Endpoints[0].Host, cConf.Endpoints[0].Port))
		if err != nil {
			return
		}
		_ = conn.Close()

		InitRedisCluster(cConf)

		assert.Equal(t, true, GetClusterInstance() != nil)

		_, err = GetClusterInstance().Ping(context.Background()).Result()
		assert.NoError(t, err)

		_, err = GetClusterInstance().Set(context.Background(), "tester-InitRedisCluster-set", "go-sail", time.Minute).Result()
		assert.NoError(t, err)

		result, err := GetClusterInstance().Get(context.Background(), "tester-InitRedisCluster-set").Result()
		assert.NoError(t, err)
		t.Log(result)
		assert.Equal(t, "go-sail", result)

		_ = GetClusterInstance().Close()
	})
}

func TestNewRedisCluster(t *testing.T) {
	t.Run("NewRedisCluster", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cConf.Endpoints[0].Host, cConf.Endpoints[0].Port))
		if err != nil {
			return
		}
		_ = conn.Close()

		newClient, err := NewCluster(cConf)
		assert.NoError(t, err)
		assert.Equal(t, true, newClient != nil)

		_, err = newClient.Ping(context.Background()).Result()
		assert.NoError(t, err)

		_, err = newClient.Set(context.Background(), "tester-NewRedisCluster-set", "go-sail", time.Minute).Result()
		assert.NoError(t, err)

		result, err := newClient.Get(context.Background(), "tester-NewRedisCluster-set").Result()
		assert.NoError(t, err)
		t.Log(result)
		assert.Equal(t, "go-sail", result)

		_ = newClient.Close()
	})
}
