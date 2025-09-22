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

func TestRedisLockPanic(t *testing.T) {
	t.Run("RedisLockPanic-Init-Panic", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				t.Log(RedisLocker())
			})
		}
	})

	t.Run("RedisLockPanic-Lock", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				key := "go-sail-redisLocker-Lock"
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				RedisLocker().Lock(ctx, key)
			})
		}
	})

	t.Run("RedisLockPanic-TryLock", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				key := "go-sail-redisLocker-TryLock"
				RedisLocker().TryLock(key)
			})
		}
	})

	t.Run("RedisLockPanic-TryLockWithContext", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				key := "go-sail-redisLocker-TryLockWithContext"
				RedisLocker().TryLock(key)
			})
		}
	})

	t.Run("RedisLockPanic-UnlockWithContext", func(t *testing.T) {
		if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
			assert.Panics(t, func() {
				key := "go-sail-redisLocker-UnlockWithContext"
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				RedisLocker().UnlockWithContext(ctx, key)
			})
		}
	})
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

	t.Run("Unlock-WithContext", func(t *testing.T) {
		key := "go-sail-redisLocker-UnlockWithContext"
		t.Log(RedisLocker().TryLock(key))
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		RedisLocker().UnlockWithContext(ctx, key)
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
			//{Host: "127.0.0.1", Port: 7003},
			//{Host: "127.0.0.1", Port: 7004},
			//{Host: "127.0.0.1", Port: 7005},
		},
	}
	//try connect
	redisClient, err := redis.NewCluster(conf)
	if err != nil || redisClient == nil {
		t.Log("redis instance not ready, this test case ignore")
		return
	}

	t.Run("InitWithClient", func(t *testing.T) {
		key := "go-sail-redisLocker-InitWithClient"
		t.Log(RedisLocker(redisClient).TryLock(key))
		assert.Equal(t, false, RedisLocker(redisClient).TryLock(key))
	})

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

	t.Run("Unlock-WithContext", func(t *testing.T) {
		key := "go-sail-redisLocker-UnlockWithContext"
		t.Log(RedisLocker().TryLock(key))
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		RedisLocker().UnlockWithContext(ctx, key)
		assert.Equal(t, true, RedisLocker().TryLock(key))
	})
}

func TestWithRedisExecuteTimeout(t *testing.T) {
	t.Run("withRedisExecuteTimeout", func(t *testing.T) {
		t.Log(withRedisExecuteTimeout())
	})
}

func TestStartRenewalScheduler(t *testing.T) {
	t.Run("startRenewalScheduler", func(t *testing.T) {
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
		rl := &redisLockerImpl{client: redisClient}
		rl.startRenewalScheduler()
		time.Sleep(time.Second * 2)
		_ = redisClient.Close()
	})
}
