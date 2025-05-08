package redis

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cConf = ClusterConf{
	Enable: true,
	Endpoints: []Endpoint{
		{Host: "127.0.0.1", Port: 7000},
		{Host: "127.0.0.1", Port: 7001},
		{Host: "127.0.0.1", Port: 7002},
		{Host: "127.0.0.1", Port: 7003},
		{Host: "127.0.0.1", Port: 7004},
		{Host: "127.0.0.1", Port: 7005},
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
	})

	t.Run("InitRedisCluster-SSLEnable", func(t *testing.T) {
		assert.Panics(t, func() {
			cConf.SSLEnable = true
			InitRedisCluster(cConf)
		})
	})
}

func TestNewRedisCluster(t *testing.T) {
	t.Run("NewRedisCluster", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cConf.Endpoints[0].Host, cConf.Endpoints[0].Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		t.Log(NewCluster(cConf))

	})

	t.Run("NewRedisCluster-SSLEnable", func(t *testing.T) {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cConf.Endpoints[0].Host, cConf.Endpoints[0].Port))
		if err != nil {
			return
		}
		_ = conn.Close()
		cConf.SSLEnable = true
		t.Log(NewCluster(cConf))
	})
}
