package utils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

type stateListenersDeprecated struct {
	mux       *sync.Mutex
	listeners map[string]chan struct{}
}

var (
	lockTTLDeprecated              = time.Second * 10
	redisExecuteTimeoutDeprecated  = time.Second * 3
	renewalCheckIntervalDeprecated = time.Second * 1
	statesDeprecated               = &stateListenersDeprecated{mux: &sync.Mutex{}, listeners: make(map[string]chan struct{})}
)

// RedisTryLock redis锁-尝试上锁（自动推测连接类型）
//
// Deprecated: RedisTryLock is deprecated,it will be removed in the future.
//
// Please use RedisLocker().TryLock() instead.
//
// using SetNX
//
// 与之对应的是使用 RedisUnlock 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func RedisTryLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeoutDeprecated)
	defer cancel()

	return RedisTryLockWithContext(ctx, key)
}

// RedisTryLockWithContext redis锁-尝试上锁（自动推测连接类型）
//
// Deprecated: RedisTryLockWithContext is deprecated,it will be removed in the future.
//
// Please use RedisLocker().TryLockWithContext() instead.
//
// using SetNX
//
// 与之对应的是使用 RedisUnlock 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func RedisTryLockWithContext(ctx context.Context, key string) bool {
	if redis.GetInstance() != nil {
		return RedisStandaloneLockWithContext(ctx, key)
	}

	if redis.GetClusterInstance() != nil {
		return RedisClusterLockWithContext(ctx, key)
	}

	panic("using redis lock on nil redis instance")
}

// RedisLock redis锁-上锁（自动推测连接类型）
//
// Deprecated: RedisLock is deprecated,it will be removed in the future.
//
// Please use RedisLocker().Lock() instead.
//
// using SetNX
//
// 与之对应的是使用 RedisUnlock 解锁
//
// # Note
//
// 该方法会阻塞住线程直到上锁成功 或者 触发ctx.Done()
func RedisLock(ctx context.Context, key string) {
	if redis.GetInstance() == nil && redis.GetClusterInstance() == nil {
		panic("using redis lock on nil redis instance")
	}

	var (
		locked = false
		ticker = time.NewTicker(time.Millisecond)
	)

LOOP:
	for {
		select {
		case <-ticker.C:
			if redis.GetInstance() != nil {
				locked = RedisStandaloneLock(key)
			}
			if redis.GetClusterInstance() != nil {
				locked = RedisClusterLock(key)
			}
			if locked {
				break LOOP
			}
		case <-ctx.Done():
			break LOOP
		}
	}

	ticker.Stop()
}

// RedisUnlock redis锁-解锁（自动推测连接类型）
//
// Deprecated: RedisUnlock is deprecated,it will be removed in the future.
//
// Please use RedisLocker().Unlock() instead.
//
// using SetNX
func RedisUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeoutDeprecated)
	go func() {
		for range time.After(redisExecuteTimeoutDeprecated) {
			cancel()
			break
		}
	}()

	RedisUnlockWithContext(ctx, key)
}

// RedisUnlockWithContext redis锁-解锁（自动推测连接类型）
//
// Deprecated: RedisUnlockWithContext is deprecated,it will be removed in the future.
//
// Please use RedisLocker().UnlockWithContext() instead.
//
// using SetNX
func RedisUnlockWithContext(ctx context.Context, key string) {
	if redis.GetInstance() != nil {
		RedisStandaloneUnlockWithContext(ctx, key)
		return
	}

	if redis.GetClusterInstance() != nil {
		RedisClusterUnlockWithContext(ctx, key)
		return
	}

	panic("using redis unlock on nil redis instance")
}

// RedisStandaloneLock redis锁-上锁（使用standalone）
//
// using SetNX
func RedisStandaloneLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeoutDeprecated)
	defer cancel()

	return RedisStandaloneLockWithContext(ctx, key)
}

// RedisStandaloneLockWithContext redis锁-上锁（使用standalone）
//
// using SetNX
func RedisStandaloneLockWithContext(ctx context.Context, key string) bool {
	ok, err := redis.GetInstance().SetNX(ctx, key, lockerValueDeprecated(), lockTTLDeprecated).Result()
	if err != nil {
		return false
	}

	if ok {
		cancelChan := make(chan struct{})

		//自动续期
		go func() {
			ticker := time.NewTicker(renewalCheckIntervalDeprecated)
			innerCtx := context.Background()
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case <-ticker.C:
					if ok, redisErr := redis.GetInstance().Expire(innerCtx, key, lockTTLDeprecated).Result(); !ok || redisErr != nil {
						break LOOP
					}
				case <-cancelChan:
					break LOOP
				}
			}
		}()

		statesDeprecated.mux.Lock()
		statesDeprecated.listeners[key] = cancelChan
		statesDeprecated.mux.Unlock()
	}

	return ok
}

// RedisStandaloneUnlock redis锁-解锁（使用standalone）
//
// using SetNX
func RedisStandaloneUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeoutDeprecated)
	defer cancel()

	RedisStandaloneUnlockWithContext(ctx, key)
}

// RedisStandaloneUnlockWithContext redis锁-解锁（使用standalone）
//
// using SetNX
func RedisStandaloneUnlockWithContext(ctx context.Context, key string) {
	_, _ = redis.GetInstance().Del(ctx, key).Result()

	go func() {
		statesDeprecated.mux.Lock()
		ch, ok := statesDeprecated.listeners[key]
		if ok {
			delete(statesDeprecated.listeners, key)
		}
		statesDeprecated.mux.Unlock()
		if ok {
			ch <- struct{}{}
			close(ch)
		}
	}()
}

// RedisClusterLock redis锁-上锁（使用cluster）
//
// using SetNX
func RedisClusterLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeoutDeprecated)
	defer cancel()

	return RedisClusterLockWithContext(ctx, key)
}

// RedisClusterLockWithContext redis锁-上锁（使用cluster）
//
// using SetNX
func RedisClusterLockWithContext(ctx context.Context, key string) bool {
	ok, err := redis.GetClusterInstance().SetNX(ctx, key, lockerValueDeprecated(), lockTTLDeprecated).Result()
	if err != nil {
		return false
	}

	if ok {
		cancelChan := make(chan struct{})

		//自动续期
		go func() {
			ticker := time.NewTicker(renewalCheckIntervalDeprecated)
			innerCtx := context.Background()
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case <-ticker.C:
					if ok, redisErr := redis.GetClusterInstance().Expire(innerCtx, key, lockTTLDeprecated).Result(); !ok || redisErr != nil {
						break LOOP
					}
				case <-cancelChan:
					break LOOP
				}
			}
		}()

		statesDeprecated.mux.Lock()
		statesDeprecated.listeners[key] = cancelChan
		statesDeprecated.mux.Unlock()
	}

	return ok
}

// RedisClusterUnlock redis锁-解锁（使用cluster）
//
// using SetNX
func RedisClusterUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeoutDeprecated)
	defer cancel()

	RedisClusterUnlockWithContext(ctx, key)
}

// RedisClusterUnlockWithContext redis锁-解锁（使用cluster）
//
// using SetNX
func RedisClusterUnlockWithContext(ctx context.Context, key string) {
	_, _ = redis.GetClusterInstance().Del(ctx, key).Result()

	go func() {
		statesDeprecated.mux.Lock()
		ch, ok := statesDeprecated.listeners[key]
		if ok {
			delete(statesDeprecated.listeners, key)
		}
		statesDeprecated.mux.Unlock()
		if ok {
			ch <- struct{}{}
			close(ch)
		}
	}()
}

// 锁的持有者信息
func lockerValueDeprecated() string {
	hostname, _ := os.Hostname()
	ip, _ := GetLocalIP()

	return fmt.Sprintf("lockedAt:%s@%s(%s)", time.Now().Format("2006-01-02T15:04:05Z"), hostname, ip)
}
