package utils

import (
	"context"
	"time"

	"github.com/keepchen/go-sail/v2/pkg/lib/redis"
)

var (
	redisExecuteTimeout  = time.Second * 3
	lockTTL              = time.Second * 10
	renewalCheckInterval = time.Second * 3
)

// RedisLock redis锁-上锁
//
// using SetNX
func RedisLock(appName, key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	//自动续期
	go func(appName, key string) {
		ticker := time.NewTicker(renewalCheckInterval)
		innerCtx := context.Background()
		defer ticker.Stop()

		for range ticker.C {
			_, err := redis.GetInstance().Get(innerCtx, WrapRedisKey(appName, key)).Result()
			if err != nil {
				ticker.Stop()
				break
			}
			_, _ = redis.GetInstance().Expire(innerCtx, WrapRedisKey(appName, key), lockTTL).Result()
		}
	}(appName, key)

	ok, err := redis.GetInstance().SetNX(ctx, WrapRedisKey(appName, key), time.Now().Unix(), lockTTL).Result()
	if err != nil {
		return false
	}

	return ok
}

// RedisUnlock redis锁-解锁
//
// using SetNX
func RedisUnlock(appName, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	_, _ = redis.GetInstance().Del(ctx, WrapRedisKey(appName, key)).Result()
}

// RedisClusterLock redis锁-上锁（使用cluster）
//
// using SetNX
func RedisClusterLock(appName, key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	//自动续期
	go func(appName, key string) {
		ticker := time.NewTicker(renewalCheckInterval)
		innerCtx := context.Background()
		defer ticker.Stop()

		for range ticker.C {
			_, err := redis.GetClusterInstance().Get(innerCtx, WrapRedisKey(appName, key)).Result()
			if err != nil {
				ticker.Stop()
				break
			}
			_, _ = redis.GetClusterInstance().Expire(innerCtx, WrapRedisKey(appName, key), lockTTL).Result()
		}
	}(appName, key)

	ok, err := redis.GetClusterInstance().SetNX(ctx, WrapRedisKey(appName, key), time.Now().Unix(), lockTTL).Result()
	if err != nil {
		return false
	}

	return ok
}

// RedisClusterUnlock redis锁-解锁（使用cluster）
//
// using SetNX
func RedisClusterUnlock(appName, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	_, _ = redis.GetClusterInstance().Del(ctx, WrapRedisKey(appName, key)).Result()
}
