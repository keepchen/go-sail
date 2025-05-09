package utils

import (
	"context"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/stretchr/testify/assert"
)

func TestLockerValue(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(lockerValueDeprecated())
	}
}

var (
	cConf = redis.ClusterConf{
		Enable: true,
		Endpoints: []redis.Endpoint{
			{Host: "127.0.0.1", Port: 7000},
			{Host: "127.0.0.1", Port: 7001},
			{Host: "127.0.0.1", Port: 7002},
			//{Host: "127.0.0.1", Port: 7003},
			//{Host: "127.0.0.1", Port: 7004},
			//{Host: "127.0.0.1", Port: 7005},
		},
	}
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
)

func TestRedisLockDeprecated(t *testing.T) {
	t.Run("TryLock-Standalone", func(t *testing.T) {
		//try connect
		redisClient, err := redis.New(sConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedis(sConf)
		key := "go-sail-redisLocker-TryLock"
		t.Log(RedisTryLock(key))
		assert.Equal(t, false, RedisTryLock(key))
	})

	t.Run("TryLock-Cluster", func(t *testing.T) {
		//try connect
		redisClient, err := redis.NewCluster(cConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedisCluster(cConf)
		key := "go-sail-redisLocker-TryLock"
		t.Log(RedisTryLock(key))
		assert.Equal(t, false, RedisTryLock(key))
	})

	t.Run("Lock-Standalone", func(t *testing.T) {
		//try connect
		redisClient, err := redis.New(sConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedis(sConf)
		key := "go-sail-redisLocker-Lock"
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		RedisLock(ctx, key)
	})

	t.Run("Lock-Cluster", func(t *testing.T) {
		//try connect
		redisClient, err := redis.NewCluster(cConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedisCluster(cConf)
		key := "go-sail-redisLocker-Lock"
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		RedisLock(ctx, key)
	})

	t.Run("Unlock-Standalone", func(t *testing.T) {
		//try connect
		redisClient, err := redis.New(sConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedis(sConf)
		key := "go-sail-redisLocker-Unlock"
		t.Log(RedisTryLock(key))
		RedisUnlock(key)
		assert.Equal(t, true, RedisTryLock(key))
	})

	t.Run("UnlockWithContext-Standalone", func(t *testing.T) {
		//try connect
		redisClient, err := redis.New(sConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedis(sConf)
		key := "go-sail-redisLocker-Unlock"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		t.Log(RedisStandaloneLockWithContext(ctx, key))
		RedisStandaloneLockWithContext(ctx, key)
	})

	t.Run("Unlock-Cluster", func(t *testing.T) {
		//try connect
		redisClient, err := redis.NewCluster(cConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedisCluster(cConf)
		key := "go-sail-redisLocker-Unlock"
		t.Log(RedisTryLock(key))
		RedisUnlock(key)
		assert.Equal(t, true, RedisTryLock(key))
	})

	t.Run("UnlockWithContext-Cluster", func(t *testing.T) {
		//try connect
		redisClient, err := redis.NewCluster(cConf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		_ = redisClient.Close()
		//initialize
		redis.InitRedisCluster(cConf)
		key := "go-sail-redisLocker-Unlock"
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		t.Log(RedisClusterLockWithContext(ctx, key))
		RedisClusterLockWithContext(ctx, key)
	})
}
