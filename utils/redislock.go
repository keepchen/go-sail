package utils

import (
	"context"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

type stateListeners struct {
	mux       *sync.Mutex
	listeners map[string]chan struct{}
}

var (
	lockTTL              = time.Second * 10
	redisExecuteTimeout  = time.Second * 3
	renewalCheckInterval = time.Second * 1
	states               = &stateListeners{mux: &sync.Mutex{}, listeners: make(map[string]chan struct{})}
)

type CancelFunc func()

// RedisTryLock redis锁-尝试上锁（自动推测连接类型）
//
// using SetNX
//
// 与之对应的是使用 RedisUnlock 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func RedisTryLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	return RedisTryLockWithContext(ctx, key)
}

// RedisTryLockWithContext redis锁-尝试上锁（自动推测连接类型）
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
// using SetNX
func RedisUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	RedisUnlockWithContext(ctx, key)
}

// RedisUnlockWithContext redis锁-解锁（自动推测连接类型）
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
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	return RedisStandaloneLockWithContext(ctx, key)
}

// RedisStandaloneLockWithContext redis锁-上锁（使用standalone）
//
// using SetNX
func RedisStandaloneLockWithContext(ctx context.Context, key string) bool {
	ok, err := redis.GetInstance().SetNX(ctx, key, time.Now().Unix(), lockTTL).Result()
	if err != nil {
		return false
	}

	if ok {
		cancelChan := make(chan struct{})

		//自动续期
		go func() {
			ticker := time.NewTicker(renewalCheckInterval)
			innerCtx := context.Background()
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case <-ticker.C:
					if ok, redisErr := redis.GetInstance().Expire(innerCtx, key, lockTTL).Result(); !ok || redisErr != nil {
						break LOOP
					}
				case <-cancelChan:
					break LOOP
				}
			}
		}()

		states.mux.Lock()
		states.listeners[key] = cancelChan
		states.mux.Unlock()
	}

	return ok
}

// RedisStandaloneUnlock redis锁-解锁（使用standalone）
//
// using SetNX
func RedisStandaloneUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	RedisStandaloneUnlockWithContext(ctx, key)
}

// RedisStandaloneUnlockWithContext redis锁-解锁（使用standalone）
//
// using SetNX
func RedisStandaloneUnlockWithContext(ctx context.Context, key string) {
	_, _ = redis.GetInstance().Del(ctx, key).Result()

	go func() {
		states.mux.Lock()
		ch, ok := states.listeners[key]
		if ok {
			delete(states.listeners, key)
		}
		states.mux.Unlock()
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
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	return RedisClusterLockWithContext(ctx, key)
}

// RedisClusterLockWithContext redis锁-上锁（使用cluster）
//
// using SetNX
func RedisClusterLockWithContext(ctx context.Context, key string) bool {
	ok, err := redis.GetClusterInstance().SetNX(ctx, key, time.Now().Unix(), lockTTL).Result()
	if err != nil {
		return false
	}

	if ok {
		cancelChan := make(chan struct{})

		//自动续期
		go func() {
			ticker := time.NewTicker(renewalCheckInterval)
			innerCtx := context.Background()
			defer ticker.Stop()

		LOOP:
			for {
				select {
				case <-ticker.C:
					if ok, redisErr := redis.GetClusterInstance().Expire(innerCtx, key, lockTTL).Result(); !ok || redisErr != nil {
						break LOOP
					}
				case <-cancelChan:
					break LOOP
				}
			}
		}()

		states.mux.Lock()
		states.listeners[key] = cancelChan
		states.mux.Unlock()
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

	RedisClusterUnlockWithContext(ctx, key)
}

// RedisClusterUnlockWithContext redis锁-解锁（使用cluster）
//
// using SetNX
func RedisClusterUnlockWithContext(ctx context.Context, key string) {
	_, _ = redis.GetClusterInstance().Del(ctx, key).Result()

	go func() {
		states.mux.Lock()
		ch, ok := states.listeners[key]
		if ok {
			delete(states.listeners, key)
		}
		states.mux.Unlock()
		if ok {
			ch <- struct{}{}
			close(ch)
		}
	}()
}
