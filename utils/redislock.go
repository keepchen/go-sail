package utils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

type redisLockerImpl struct {
}

type IRedisLocker interface {
	// TryLock redis锁-尝试上锁（自动推测连接类型）
	//
	// using SetNX
	//
	// 与之对应的是使用 Unlock 解锁
	//
	// # Note
	//
	// 该方法会立即返回锁定成功与否的结果
	TryLock(key string) bool
	// TryLockWithContext redis锁-尝试上锁（自动推测连接类型）
	//
	// using SetNX
	//
	// 与之对应的是使用 Unlock 解锁
	//
	// # Note
	//
	// 该方法会立即返回锁定成功与否的结果
	TryLockWithContext(ctx context.Context, key string) bool
	// Lock redis锁-上锁（自动推测连接类型）
	//
	// using SetNX
	//
	// 与之对应的是使用 Unlock 解锁
	//
	// # Note
	//
	// 该方法会阻塞住线程直到上锁成功 或者 触发ctx.Done()
	Lock(ctx context.Context, key string)
	// Unlock redis锁-解锁（自动推测连接类型）
	//
	// using SetNX
	Unlock(key string)
	// UnlockWithContext redis锁-解锁（自动推测连接类型）
	//
	// using SetNX
	UnlockWithContext(ctx context.Context, key string)
}

var _ IRedisLocker = &redisLockerImpl{}

// RedisLocker 实例化redis锁工具类
func RedisLocker() IRedisLocker {
	return &redisLockerImpl{}
}

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

// TryLock redis锁-尝试上锁（自动推测连接类型）
//
// using SetNX
//
// 与之对应的是使用 Unlock 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func (rl redisLockerImpl) TryLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	return rl.TryLockWithContext(ctx, key)
}

// TryLockWithContext redis锁-尝试上锁（自动推测连接类型）
//
// using SetNX
//
// 与之对应的是使用 Unlock 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func (rl redisLockerImpl) TryLockWithContext(ctx context.Context, key string) bool {
	if redis.GetInstance() != nil {
		return rl.StandaloneLockWithContext(ctx, key)
	}

	if redis.GetClusterInstance() != nil {
		return rl.ClusterLockWithContext(ctx, key)
	}

	panic("using redis lock on nil redis instance")
}

// Lock redis锁-上锁（自动推测连接类型）
//
// using SetNX
//
// 与之对应的是使用 Unlock 解锁
//
// # Note
//
// 该方法会阻塞住线程直到上锁成功 或者 触发ctx.Done()
func (rl redisLockerImpl) Lock(ctx context.Context, key string) {
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
				locked = rl.StandaloneLock(key)
			}
			if redis.GetClusterInstance() != nil {
				locked = rl.ClusterLock(key)
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

// Unlock redis锁-解锁（自动推测连接类型）
//
// using SetNX
func (rl redisLockerImpl) Unlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	go func() {
		for range time.After(redisExecuteTimeout) {
			cancel()
			break
		}
	}()

	rl.UnlockWithContext(ctx, key)
}

// UnlockWithContext redis锁-解锁（自动推测连接类型）
//
// using SetNX
func (rl redisLockerImpl) UnlockWithContext(ctx context.Context, key string) {
	if redis.GetInstance() != nil {
		rl.StandaloneUnlockWithContext(ctx, key)
		return
	}

	if redis.GetClusterInstance() != nil {
		rl.ClusterUnlockWithContext(ctx, key)
		return
	}

	panic("using redis unlock on nil redis instance")
}

// StandaloneLock redis锁-上锁（使用standalone）
//
// using SetNX
func (rl redisLockerImpl) StandaloneLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	return rl.StandaloneLockWithContext(ctx, key)
}

// StandaloneLockWithContext redis锁-上锁（使用standalone）
//
// using SetNX
func (rl redisLockerImpl) StandaloneLockWithContext(ctx context.Context, key string) bool {
	ok, err := redis.GetInstance().SetNX(ctx, key, lockerValue(), lockTTL).Result()
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

// StandaloneUnlock redis锁-解锁（使用standalone）
//
// using SetNX
func (rl redisLockerImpl) StandaloneUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	rl.StandaloneUnlockWithContext(ctx, key)
}

// StandaloneUnlockWithContext redis锁-解锁（使用standalone）
//
// using SetNX
func (redisLockerImpl) StandaloneUnlockWithContext(ctx context.Context, key string) {
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

// ClusterLock redis锁-上锁（使用cluster）
//
// using SetNX
func (rl redisLockerImpl) ClusterLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	return rl.ClusterLockWithContext(ctx, key)
}

// ClusterLockWithContext redis锁-上锁（使用cluster）
//
// using SetNX
func (redisLockerImpl) ClusterLockWithContext(ctx context.Context, key string) bool {
	ok, err := redis.GetClusterInstance().SetNX(ctx, key, lockerValue(), lockTTL).Result()
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

// ClusterUnlock redis锁-解锁（使用cluster）
//
// using SetNX
func (rl redisLockerImpl) ClusterUnlock(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	rl.ClusterUnlockWithContext(ctx, key)
}

// ClusterUnlockWithContext redis锁-解锁（使用cluster）
//
// using SetNX
func (redisLockerImpl) ClusterUnlockWithContext(ctx context.Context, key string) {
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

// 锁的持有者信息
func lockerValue() string {
	hostname, _ := os.Hostname()
	ip, _ := IP().GetLocal()

	return fmt.Sprintf("lockedAt:%s@%s(%s)", time.Now().Format("2006-01-02T15:04:05Z"), hostname, ip)
}
