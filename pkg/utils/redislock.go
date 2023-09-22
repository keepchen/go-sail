package utils

import (
	"context"
	"time"

	"github.com/keepchen/go-sail/v2/pkg/lib/redis"
)

var (
	lockTTL                         = time.Second * 10
	redisExecuteTimeout             = time.Second * 3
	renewalCheckInterval            = time.Second * 3
	cancelRenewalFuncChannel        = make(chan struct{})
	cancelRenewalFuncChannelCluster = make(chan struct{})
)

// RedisLock redis锁-上锁
//
// using SetNX
func RedisLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	ok, err := redis.GetInstance().SetNX(ctx, key, time.Now().Unix(), lockTTL).Result()
	if err != nil {
		return false
	}

	if ok {
		//自动续期
		go func(key string) {
			ticker := time.NewTicker(renewalCheckInterval)
			innerCtx := context.Background()
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case <-ticker.C:
					_, redisErr := redis.GetInstance().Get(innerCtx, key).Result()
					if redisErr != nil {
						break LOOP
					}
					_, _ = redis.GetInstance().Expire(innerCtx, key, lockTTL).Result()
				case <-cancelRenewalFuncChannel:
					break LOOP
				}
			}
		}(key)
	}

	return ok
}

// RedisUnlock redis锁-解锁
//
// using SetNX
func RedisUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	_, _ = redis.GetInstance().Del(ctx, key).Result()

	go func() {
		cancelRenewalFuncChannel <- struct{}{}
	}()
}

// RedisClusterLock redis锁-上锁（使用cluster）
//
// using SetNX
func RedisClusterLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	ok, err := redis.GetClusterInstance().SetNX(ctx, key, time.Now().Unix(), lockTTL).Result()
	if err != nil {
		return false
	}

	if ok {
		//自动续期
		go func(key string) {
			ticker := time.NewTicker(renewalCheckInterval)
			innerCtx := context.Background()
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case <-ticker.C:
					_, redisErr := redis.GetClusterInstance().Get(innerCtx, key).Result()
					if redisErr != nil {
						break LOOP
					}
					_, _ = redis.GetClusterInstance().Expire(innerCtx, key, lockTTL).Result()
				case <-cancelRenewalFuncChannelCluster:
					break LOOP
				}
			}
		}(key)
	}

	return ok
}

// RedisClusterUnlock redis锁-解锁（使用cluster）
//
// using SetNX
func RedisClusterUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	_, _ = redis.GetClusterInstance().Del(ctx, key).Result()

	go func() {
		cancelRenewalFuncChannelCluster <- struct{}{}
	}()
}
