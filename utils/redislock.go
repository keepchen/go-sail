package utils

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	redisLib "github.com/go-redis/redis/v8"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

type redisLockerImpl struct {
	client redisLib.UniversalClient
}

// IRedisLocker redis锁定义
type IRedisLocker interface {
	// TryLock redis锁-尝试上锁
	//
	// using SetNX
	//
	// 与之对应的是使用 Unlock 或 UnlockWithContext 解锁
	//
	// # Note
	//
	// 该方法会立即返回锁定成功与否的结果
	TryLock(key string) bool
	// TryLockWithContext redis锁-尝试上锁
	//
	// using SetNX
	//
	// 与之对应的是使用 Unlock 或 UnlockWithContext 解锁
	//
	// # Note
	//
	// 该方法会立即返回锁定成功与否的结果
	TryLockWithContext(ctx context.Context, key string) bool
	// Lock redis锁-上锁
	//
	// using SetNX
	//
	// 与之对应的是使用 Unlock 或 UnlockWithContext 解锁
	//
	// # Note
	//
	// 该方法会阻塞住线程直到上锁成功 或者 触发ctx.Done()
	Lock(ctx context.Context, key string)
	// Unlock redis锁-解锁
	//
	// using Del
	//
	// # Note
	//
	// 该方法会立即返回解锁成功与否的结果
	Unlock(key string) bool
	// UnlockWithContext redis锁-解锁
	//
	// using Del
	//
	// # Note
	//
	// 该方法会阻塞住线程直到解锁成功 或者 触发ctx.Done()
	UnlockWithContext(ctx context.Context, key string)
}

var _ IRedisLocker = &redisLockerImpl{}

var redisLockerPool = sync.Pool{
	New: func() interface{} {
		return &redisLockerImpl{}
	},
}

// RedisLocker 实例化redis锁工具类
//
// # Note
//
// 若未使用自定义客户端且单实例和集群客户端都没有实例化，那么将panic
func RedisLocker(client ...redisLib.UniversalClient) IRedisLocker {
	rl := redisLockerPool.Get().(*redisLockerImpl)
	//使用自定义客户端
	if len(client) > 0 {
		rl.client = client[0]
		return rl
	}
	//使用单实例客户端
	if redis.GetInstance() != nil {
		rl.client = redis.GetInstance()
		return rl
	}
	//使用集群客户端
	if redis.GetClusterInstance() != nil {
		rl.client = redis.GetClusterInstance()
		return rl
	}
	redisLockerPool.Put(rl)
	panic("using redis lock on nil redis instance")
}

type stateListeners struct {
	mux       *sync.Mutex
	listeners map[string]chan struct{}
}

var (
	lockTTL              = time.Second * 10
	redisExecuteTimeout  = time.Second * 3
	retryInterval        = time.Millisecond * 100
	renewalCheckInterval = time.Second * 1
	states               = &stateListeners{mux: &sync.Mutex{}, listeners: make(map[string]chan struct{})}
)

// TryLock redis锁-尝试上锁
//
// using SetNX
//
// 与之对应的是使用 Unlock 或 UnlockWithContext 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func (rl *redisLockerImpl) TryLock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	return rl.TryLockWithContext(ctx, key)
}

// TryLockWithContext redis锁-尝试上锁
//
// using SetNX
//
// 与之对应的是使用 Unlock 或 UnlockWithContext 解锁
//
// # Note
//
// 该方法会立即返回锁定成功与否的结果
func (rl *redisLockerImpl) TryLockWithContext(ctx context.Context, key string) bool {
	lockOk, _ := rl.client.SetNX(ctx, key, lockerValue(), lockTTL).Result()

	//锁定成功，开始执行自动续期
	if lockOk {
		rl.autoRenewal(key)
	}

	redisLockerPool.Put(rl)

	return lockOk
}

// Lock redis锁-上锁
//
// using SetNX
//
// 与之对应的是使用 Unlock 或 UnlockWithContext 解锁
//
// # Note
//
// 该方法会阻塞住线程直到上锁成功 或者 触发ctx.Done()
func (rl *redisLockerImpl) Lock(ctx context.Context, key string) {
	lockOk, lockErr := rl.client.SetNX(ctx, key, lockerValue(), lockTTL).Result()

	//第一次锁定失败，进行重试操作
	if !lockOk || lockErr != nil {
		retryTicker := time.NewTicker(retryInterval)

	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-retryTicker.C:
				lockOk, lockErr = rl.client.SetNX(ctx, key, lockerValue(), lockTTL).Result()
				if lockOk && lockErr == nil {
					break LOOP
				}
			}
		}
	}

	//锁定成功，开始执行自动续期
	if lockOk {
		rl.autoRenewal(key)
	}

	redisLockerPool.Put(rl)
}

// Unlock redis锁-解锁
//
// using Del
func (rl *redisLockerImpl) Unlock(key string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), redisExecuteTimeout)
	defer cancel()

	unlockOk, _ := rl.client.Del(ctx, key).Result()

	//清理内存数据并终止自动续期
	//
	// 调用解锁方法意图明显，因此无论是否解锁成功，都执行收尾工作
	rl.clearListenerAndStopAutoRenewal(key)

	redisLockerPool.Put(rl)

	return unlockOk == 1
}

// UnlockWithContext redis锁-解锁
//
// using Del
func (rl *redisLockerImpl) UnlockWithContext(ctx context.Context, key string) {
	unlockOk, unlockErr := rl.client.Del(ctx, key).Result()

	if unlockOk != 1 || unlockErr != nil {
		ticker := time.NewTicker(retryInterval)

	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-ticker.C:
				unlockOk, unlockErr = rl.client.Del(ctx, key).Result()
				if unlockOk > 0 && unlockErr == nil {
					break LOOP
				}
			}
		}
	}

	//清理内存数据并终止自动续期
	rl.clearListenerAndStopAutoRenewal(key)

	redisLockerPool.Put(rl)
}

// 自动续期
func (rl *redisLockerImpl) autoRenewal(key string) {
	cancelChan := make(chan struct{})

	go func() {
		ticker := time.NewTicker(renewalCheckInterval)
		innerCtx := context.Background()
		defer ticker.Stop()

	LOOP:
		for {
			select {
			case <-ticker.C:
				if expOk, expErr := rl.client.ExpireXX(innerCtx, key, lockTTL).Result(); !expOk || expErr != nil {
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

// 清理监听器并停止自动续期
func (rl *redisLockerImpl) clearListenerAndStopAutoRenewal(key string) {
	states.mux.Lock()
	ch, ok := states.listeners[key]
	if ok {
		delete(states.listeners, key)
	}
	states.mux.Unlock()
	if ok {
		go func() {
			ch <- struct{}{}
			close(ch)
		}()
	}
}

// 锁的持有者信息
func lockerValue() string {
	hostname, _ := os.Hostname()
	ip, _ := IP().GetLocal()

	return fmt.Sprintf("lockedAt:%s@%s(%s)", hostname, ip, time.Now().Format("2006-01-02T15:04:05Z"))
}
