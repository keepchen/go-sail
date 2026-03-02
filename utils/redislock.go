package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
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
	// 该方法会阻塞住线程直到解锁有结果 或者 触发ctx.Done()
	UnlockWithContext(ctx context.Context, key string) bool
}

var _ IRedisLocker = &redisLockerImpl{}

var onceStarRenewalScheduler sync.Once

// RedisLocker 实例化redis锁工具类
//
// # Note
//
// 1.若未指定自定义客户端且单实例和集群客户端都没有实例化，那么将panic
//
// 2.若指定了自定义客户端，请始终保持相同的客户端调用，否则将造成数据异常
func RedisLocker(client ...redisLib.UniversalClient) IRedisLocker {
	rl := &redisLockerImpl{}

	defer onceStarRenewalScheduler.Do(rl.startRenewalScheduler)

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
	panic("using redis lock on nil redis instance")
}

type cancelControl struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type stateListeners struct {
	mux       *sync.RWMutex
	listeners map[string]*cancelControl
}

var (
	lockTTL              = time.Second * 10
	redisExecuteTimeout  = time.Second * 3
	retryInterval        = time.Millisecond * 100
	renewalCheckInterval = time.Second * 1
	states               = &stateListeners{mux: &sync.RWMutex{}, listeners: make(map[string]*cancelControl)}
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
	ctx, cancel := withRedisExecuteTimeout()
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
	if !canDoLockPreflight(key) {
		return false
	}

	lockOk, lockErr := rl.client.SetNX(ctx, key, lockerValue(), lockTTL).Result()
	if lockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s lock err: %v\n", key, lockErr)
	}

	//锁定成功，开始执行自动续期
	if lockOk {
		rl.autoRenewal(key)
	}

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

	if lockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s lock error: %v\n", key, lockErr)
	}

	//第一次锁定失败，进行重试操作
	if !lockOk || lockErr != nil {
		retryTicker := time.NewTicker(retryInterval)
		defer retryTicker.Stop()

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
}

// Unlock redis锁-解锁
//
// using Del
func (rl *redisLockerImpl) Unlock(key string) bool {
	ctx, cancel := withRedisExecuteTimeout()
	defer cancel()

	//持有者一致性检测(如果获取失败，也认为不符合一致性)
	lv, err := rl.client.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s unlock get key err: %v\n", key, err)
		return false
	}
	if !holderConsistencyDetection(lv) {
		return false
	}

	unlockOk, unlockErr := rl.client.Del(ctx, key).Result()
	if unlockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s unlock delete key error: %v\n", key, unlockErr)
	}

	//清理内存数据并终止自动续期
	//
	// 调用解锁方法意图明显，因此无论是否解锁成功，都执行收尾工作
	rl.clearListenerAndStopAutoRenewal(key)

	return unlockOk == 1
}

// UnlockWithContext redis锁-解锁
//
// using Del
func (rl *redisLockerImpl) UnlockWithContext(ctx context.Context, key string) bool {
	//持有者一致性检测(如果获取失败，也认为不符合一致性)
	lv, err := rl.client.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s unlock with context get key err: %v\n", key, err)
		return false
	}
	if !holderConsistencyDetection(lv) {
		return false
	}
	unlockOk, unlockErr := rl.client.Del(ctx, key).Result()

	if unlockOk != 1 || unlockErr != nil {
		ticker := time.NewTicker(retryInterval)
		defer ticker.Stop()

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

	return unlockOk == 1
}

// 自动续期
func (rl *redisLockerImpl) autoRenewal(key string) {
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := &cancelControl{ctx: ctx, cancel: cancel}

	states.mux.Lock()
	states.listeners[key] = ctrl
	states.mux.Unlock()
}

// 清理监听器并停止自动续期
func (rl *redisLockerImpl) clearListenerAndStopAutoRenewal(key string) {
	states.mux.Lock()
	if ctrl, ok := states.listeners[key]; ok {
		ctrl.cancel()
		delete(states.listeners, key)
	}
	states.mux.Unlock()
}

type keyAndCtrl struct {
	key  string
	ctrl *cancelControl
}

// 续期的统一调度器
//
// 1.使用二阶段Mutex锁定，减少锁持有时间
//
// 2.使用redis pipeline减少RTT
func (rl *redisLockerImpl) startRenewalScheduler() {
	doRenewalRound := func() {
		if rl.client == nil {
			fmt.Println("[Go-Sail] <redisLock> renewal task not emit cause redis client is nil")
			return
		}

		states.mux.Lock()
		if len(states.listeners) == 0 {
			states.mux.Unlock()
			//避免空转锁定占用
			return
		}
		processingKeys := make([]*keyAndCtrl, 0, len(states.listeners))
		for key, ctrl := range states.listeners {
			processingKeys = append(processingKeys, &keyAndCtrl{key: key, ctrl: ctrl})
		}
		states.mux.Unlock()

		ctx, cancel := withRedisExecuteTimeout()
		defer cancel()

		validKeys := make([]*keyAndCtrl, 0, len(processingKeys))
		invalidKeys := make([]*keyAndCtrl, 0, len(processingKeys))
		cmds, pipeErr := rl.client.Pipelined(ctx, func(pipe redisLib.Pipeliner) error {
			for index := range processingKeys {
				if processingKeys[index].ctrl.ctx.Err() != nil {
					invalidKeys = append(invalidKeys, processingKeys[index])
				} else {
					pipe.Expire(ctx, processingKeys[index].key, lockTTL)
					validKeys = append(validKeys, processingKeys[index])
				}
			}
			return nil
		})

		if pipeErr != nil {
			fmt.Printf("[Go-Sail] <redisLock> renewal pipeline error: %v\n", pipeErr)
			return //宽容处理：本轮不进行cancel和delete
		}

		states.mux.Lock()
		for index := range cmds {
			if expOk, expErr := cmds[index].(*redisLib.BoolCmd).Result(); !expOk || expErr != nil {
				if expErr != nil {
					fmt.Printf("[Go-Sail] <redisLock> key: %s renewal err: %v\n", validKeys[index].key, expErr)
				}
				validKeys[index].ctrl.cancel() //续期失败也要清理掉
				delete(states.listeners, validKeys[index].key)
			}
		}
		//清理已经过期的
		for index := range invalidKeys {
			invalidKeys[index].ctrl.cancel() //保险的再次调用以确保触发ctx.Done
			delete(states.listeners, invalidKeys[index].key)
		}
		states.mux.Unlock()
	}

	go func() {
		ticker := time.NewTicker(renewalCheckInterval)
		defer ticker.Stop()

		for range ticker.C {
			doRenewalRound()
		}
	}()
}

// 预检是否可以执行锁定任务
//
// 此操作属于本地(堆栈)检测
func canDoLockPreflight(key string) bool {
	states.mux.RLock()
	defer states.mux.RUnlock()
	_, exist := states.listeners[key]
	return !exist
}

// redis操作超时控制
func withRedisExecuteTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), redisExecuteTimeout)
}

var (
	hostname, _ = os.Hostname()   //主机名称
	ip, _       = IP().GetLocal() //主机ip
	processId   = os.Getpid()     //进程id
)

// 锁的持有者信息
func lockerValue() string {
	return fmt.Sprintf("lockedAt:%s@%s<%d>(%s)",
		hostname, ip, processId, time.Now().Format("2006-01-02T15:04:05.000000Z"))
}

// 持有者一致性检测
//
// # 注意：
//
// 一致性检测以【机器主机名+ip+进程id】为判断依据，
//
// 这样设计为的是锁只能被【持有者自己】释放，若进程
//
// down掉，堆栈中的自动维护信息会被释放，
//
// 因此即便是重新启动获取到了相同的进程号，也不受影响。
func holderConsistencyDetection(lockerValue string) bool {
	return strings.HasPrefix(lockerValue, fmt.Sprintf("lockedAt:%s@%s<%d>(", hostname, ip, processId))
}
