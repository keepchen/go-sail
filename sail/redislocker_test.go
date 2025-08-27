package sail

import (
	"fmt"
	"net"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/stretchr/testify/assert"
)

func TestSetRedisClientOnce(t *testing.T) {
	t.Run("SetRedisClientOnce", func(t *testing.T) {
		SetRedisClientOnce(nil)
	})
}

func TestRedisLocker(t *testing.T) {
	t.Run("RedisLocker-Nil-Client", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				t.Log(RedisLocker())
			})
		} else {
			t.Log(RedisLocker())
		}
	})
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", sConf.Host, sConf.Port))
	if err != nil {
		return
	}
	_ = conn.Close()
	t.Run("RedisLocker", func(t *testing.T) {
		client, err := redis.New(sConf)
		t.Log(client, err)
		SetRedisClientOnce(client)
		t.Log(RedisLocker())
	})
}
