package logger

import (
	"testing"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

var (
	sConf = redis.Conf{
		Enable: true,
		Endpoint: redis.Endpoint{
			Host:     "127.0.0.1",
			Port:     6379,
			Username: "",
			Password: "",
		},
		Database: 0,
	}
	cConf = redis.ClusterConf{
		Enable: true,
		Endpoints: []redis.Endpoint{
			{Host: "127.0.0.1", Port: 7000},
			{Host: "127.0.0.1", Port: 7001},
			{Host: "127.0.0.1", Port: 7002},
			{Host: "127.0.0.1", Port: 7003},
			{Host: "127.0.0.1", Port: 7004},
			{Host: "127.0.0.1", Port: 7005},
		},
	}
)

func TestRedisSync(t *testing.T) {
	t.Run("Sync-Standalone", func(t *testing.T) {
		writer := &redisWriterStd{}
		t.Log(writer.Sync())
	})

	t.Run("Sync-Cluster", func(t *testing.T) {
		writer := &redisClusterWriterStd{}
		t.Log(writer.Sync())
	})
}

func TestRedisWrite(t *testing.T) {
	t.Run("Write-Standalone", func(t *testing.T) {
		rd, err := redis.New(sConf)
		t.Log(err)
		if rd == nil {
			return
		}
		writer := &redisWriterStd{
			cli:     rd,
			listKey: "go-sail-tester-logger-list",
		}
		t.Log(writer.Sync())
	})

	t.Run("Write-Cluster", func(t *testing.T) {
		rd, err := redis.NewCluster(cConf)
		t.Log(err)
		if rd == nil {
			return
		}
		writer := &redisClusterWriterStd{
			cli:     rd,
			listKey: "go-sail-tester-logger-list",
		}
		t.Log(writer.Sync())
	})
}
