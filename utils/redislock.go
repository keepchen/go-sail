package utils

import (
	"context"
	"crypto/rand"
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
	token  string
}

// lua脚本-原子操作-对比并删除
const compareAndDeleteScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

// lua脚本-原子操作-对比并续期
const compareAndExpireScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0
`

var releaseLockScript = redisLib.NewScript(compareAndDeleteScript)

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

	token := lockerValue()
	lockOk, lockErr := rl.client.SetNX(ctx, key, token, lockTTL).Result()
	if lockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s lock err: %v\n", key, lockErr)
	}

	//锁定成功，开始执行自动续期
	if lockOk {
		rl.autoRenewal(key, token)
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
	token := lockerValue()
	lockOk, lockErr := rl.client.SetNX(ctx, key, token, lockTTL).Result()

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
				token = lockerValue()
				lockOk, lockErr = rl.client.SetNX(ctx, key, token, lockTTL).Result()
				if lockOk && lockErr == nil {
					break LOOP
				}
			}
		}
	}

	//锁定成功，开始执行自动续期
	if lockOk {
		rl.autoRenewal(key, token)
	}
}

// Unlock redis锁-解锁
//
// 使用Lua脚本原子校验所有权令牌并删除锁
func (rl *redisLockerImpl) Unlock(key string) bool {
	ctx, cancel := withRedisExecuteTimeout()
	defer cancel()

	ctrl, ok := listener(key)
	if !ok {
		return false
	}

	unlockOk, unlockErr := releaseLockScript.Run(ctx, rl.client, []string{key}, ctrl.token).Int64()
	if unlockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s unlock key error: %v\n", key, unlockErr)
	}

	//清理内存数据并终止自动续期
	//
	// 调用解锁方法意图明显，因此无论是否解锁成功，都执行收尾工作
	rl.clearListenerAndStopAutoRenewal(key, ctrl)

	return unlockOk == 1
}

// UnlockWithContext redis锁-解锁
//
// 使用Lua脚本原子校验所有权令牌并删除锁
func (rl *redisLockerImpl) UnlockWithContext(ctx context.Context, key string) bool {
	ctrl, ok := listener(key)
	if !ok {
		return false
	}
	unlockOk, unlockErr := releaseLockScript.Run(ctx, rl.client, []string{key}, ctrl.token).Int64()

	if unlockErr != nil {
		ticker := time.NewTicker(retryInterval)
		defer ticker.Stop()

	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-ticker.C:
				unlockOk, unlockErr = releaseLockScript.Run(ctx, rl.client, []string{key}, ctrl.token).Int64()
				if unlockErr == nil {
					break LOOP
				}
			}
		}
	}

	//清理内存数据并终止自动续期
	rl.clearListenerAndStopAutoRenewal(key, ctrl)

	return unlockOk == 1
}

// 自动续期
func (rl *redisLockerImpl) autoRenewal(key, token string) {
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := &cancelControl{ctx: ctx, cancel: cancel, token: token}

	states.mux.Lock()
	states.listeners[key] = ctrl
	states.mux.Unlock()
}

// 清理监听器并停止自动续期
func (rl *redisLockerImpl) clearListenerAndStopAutoRenewal(key string, expected *cancelControl) {
	states.mux.Lock()
	if ctrl, ok := states.listeners[key]; ok && ctrl == expected {
		ctrl.cancel()
		delete(states.listeners, key)
	}
	states.mux.Unlock()
}

type keyAndCtrl struct {
	key  string
	ctrl *cancelControl
	cmd  *redisLib.Cmd
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

		validCommands := make([]*keyAndCtrl, 0, len(processingKeys))
		invalidKeys := make([]*keyAndCtrl, 0, len(processingKeys))
		_, pipeErr := rl.client.Pipelined(ctx, func(pipe redisLib.Pipeliner) error {
			for index := range processingKeys {
				if processingKeys[index].ctrl.ctx.Err() != nil {
					invalidKeys = append(invalidKeys, processingKeys[index])
				} else {
					processingKeys[index].cmd = pipe.Eval(ctx, compareAndExpireScript, []string{processingKeys[index].key},
						processingKeys[index].ctrl.token, lockTTL.Milliseconds())
					validCommands = append(validCommands, processingKeys[index])
				}
			}
			return nil
		})

		if pipeErr != nil {
			fmt.Printf("[Go-Sail] <redisLock> renewal pipeline error: %v\n", pipeErr)
			return //宽容处理：本轮不进行cancel和delete
		}

		states.mux.Lock()
		for _, command := range validCommands {
			expResult, expErr := command.cmd.Int64()
			if expResult != 1 || expErr != nil {
				if expErr != nil {
					fmt.Printf("[Go-Sail] <redisLock> key: %s renewal err: %v\n", command.key, expErr)
				}
				if current, ok := states.listeners[command.key]; ok && current == command.ctrl {
					current.cancel() //续期失败也要清理掉
					delete(states.listeners, command.key)
				}
			}
		}
		//清理已经过期的
		for index := range invalidKeys {
			if current, ok := states.listeners[invalidKeys[index].key]; ok && current == invalidKeys[index].ctrl {
				current.cancel() //保险的再次调用以确保触发ctx.Done
				delete(states.listeners, invalidKeys[index].key)
			}
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

func listener(key string) (*cancelControl, bool) {
	states.mux.RLock()
	defer states.mux.RUnlock()
	ctrl, ok := states.listeners[key]
	return ctrl, ok
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
	return fmt.Sprintf("lockedAt:%s@%s<%d>(%s:%s)",
		hostname, ip, processId, time.Now().Format("2006-01-02T15:04:05.000000Z"), rand.Text())
}
