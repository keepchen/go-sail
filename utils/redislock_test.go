package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

func TestRedisLockerImplLockerValue(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(lockerValue())
	}
}

func TestRedisLock(t *testing.T) {
	conf := redis.Conf{
		Endpoint: redis.Endpoint{
			Host: "127.0.0.1",
			Port: 6379,
		},
	}
	//try connect
	redisClient, err := redis.New(conf)
	if err != nil || redisClient == nil {
		t.Log("redis instance not ready, this test case ignore")
		return
	}
	_ = redisClient.Close()

	//initialize
	redis.InitRedis(conf)

	t.Run("TryLock", func(t *testing.T) {
		key := "go-sail-redisLocker-TryLock"
		t.Log(RedisLocker().TryLock(key))
		assert.Equal(t, false, RedisLocker().TryLock(key))
	})

	t.Run("Lock", func(t *testing.T) {
		key := "go-sail-redisLocker-Lock"
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		RedisLocker().Lock(ctx, key)
	})

	t.Run("Unlock", func(t *testing.T) {
		key := "go-sail-redisLocker-Unlock"
		t.Log(RedisLocker().TryLock(key))
		RedisLocker().Unlock(key)
		assert.Equal(t, true, RedisLocker().TryLock(key))
	})
}

func TestRedisClusterLock(t *testing.T) {
	conf := redis.ClusterConf{
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
	//try connect
	redisClient, err := redis.NewCluster(conf)
	if err != nil || redisClient == nil {
		t.Log("redis instance not ready, this test case ignore")
		return
	}
	_ = redisClient.Close()

	//initialize
	redis.InitRedisCluster(conf)

	t.Run("TryLock", func(t *testing.T) {
		key := "go-sail-redisLocker-TryLock"
		t.Log(RedisLocker().TryLock(key))
		assert.Equal(t, false, RedisLocker().TryLock(key))
	})

	t.Run("Lock", func(t *testing.T) {
		key := "go-sail-redisLocker-Lock"
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		RedisLocker().Lock(ctx, key)
	})

	t.Run("Unlock", func(t *testing.T) {
		key := "go-sail-redisLocker-Unlock"
		t.Log(RedisLocker().TryLock(key))
		RedisLocker().Unlock(key)
		assert.Equal(t, true, RedisLocker().TryLock(key))
	})
}
