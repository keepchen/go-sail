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

// IRedisLocker Redis lock definition
type IRedisLocker interface {
	// TryLock Redis lock - try to acquire
	//
	// using SetNX
	//
	// Use Unlock or UnlockWithContext to release the lock
	//
	// # Note
	//
	// This method immediately returns whether locking succeeded
	TryLock(key string) bool
	// TryLockWithContext Redis lock - try to acquire
	//
	// using SetNX
	//
	// Use Unlock or UnlockWithContext to release the lock
	//
	// # Note
	//
	// This method immediately returns whether locking succeeded
	TryLockWithContext(ctx context.Context, key string) bool
	// Lock Redis lock - acquire
	//
	// using SetNX
	//
	// Use Unlock or UnlockWithContext to release the lock
	//
	// # Note
	//
	// This method blocks until the lock is acquired or ctx.Done() is triggered
	Lock(ctx context.Context, key string)
	// Unlock Redis lock - release
	//
	// using Del
	//
	// # Note
	//
	// This method immediately returns whether unlocking succeeded
	Unlock(key string) bool
	// UnlockWithContext Redis lock - release
	//
	// using Del
	//
	// # Note
	//
	// This method blocks until unlocking produces a result or ctx.Done() is triggered
	UnlockWithContext(ctx context.Context, key string) bool
}

var _ IRedisLocker = &redisLockerImpl{}

var onceStartRenewalScheduler sync.Once

// RedisLocker instantiates Redis lock utilities
//
// # Note
//
// 1.Panics if no custom client is specified and neither the standalone nor cluster client has been initialized
//
// 2.If a custom client is specified, consistently use the same client for calls to avoid data inconsistencies
func RedisLocker(client ...redisLib.UniversalClient) IRedisLocker {
	rl := &redisLockerImpl{}

	defer onceStartRenewalScheduler.Do(rl.startRenewalScheduler)

	//Use a custom client
	if len(client) > 0 {
		rl.client = client[0]
		return rl
	}
	//Use the standalone client
	if redis.GetInstance() != nil {
		rl.client = redis.GetInstance()
		return rl
	}
	//Use the cluster client
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

// Lua script - atomic operation - compare and delete
const compareAndDeleteScript = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`

// Lua script - atomic operation - compare and renew
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

// TryLock Redis lock - try to acquire
//
// using SetNX
//
// Use Unlock or UnlockWithContext to release the lock
//
// # Note
//
// This method immediately returns whether locking succeeded
func (rl *redisLockerImpl) TryLock(key string) bool {
	ctx, cancel := withRedisExecuteTimeout()
	defer cancel()

	return rl.TryLockWithContext(ctx, key)
}

// TryLockWithContext Redis lock - try to acquire
//
// using SetNX
//
// Use Unlock or UnlockWithContext to release the lock
//
// # Note
//
// This method immediately returns whether locking succeeded
func (rl *redisLockerImpl) TryLockWithContext(ctx context.Context, key string) bool {
	if !canDoLockPreflight(key) {
		return false
	}

	token := lockerValue()
	lockOk, lockErr := rl.client.SetNX(ctx, key, token, lockTTL).Result()
	if lockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s lock err: %v\n", key, lockErr)
	}

	//Lock acquired; start automatic renewal
	if lockOk {
		rl.autoRenewal(key, token)
	}

	return lockOk
}

// Lock Redis lock - acquire
//
// using SetNX
//
// Use Unlock or UnlockWithContext to release the lock
//
// # Note
//
// This method blocks until the lock is acquired or ctx.Done() is triggered
func (rl *redisLockerImpl) Lock(ctx context.Context, key string) {
	token := lockerValue()
	lockOk, lockErr := rl.client.SetNX(ctx, key, token, lockTTL).Result()

	if lockErr != nil {
		fmt.Printf("[Go-Sail] <redisLock> key: %s lock error: %v\n", key, lockErr)
	}

	//Initial lock attempt failed; retry
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

	//Lock acquired; start automatic renewal
	if lockOk {
		rl.autoRenewal(key, token)
	}
}

// Unlock Redis lock - release
//
// Use a Lua script to atomically verify the ownership token and delete the lock
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

	//Clean up in-memory data and stop automatic renewal
	//
	// An unlock call expresses clear intent, so perform cleanup whether unlocking succeeds or not
	rl.clearListenerAndStopAutoRenewal(key, ctrl)

	return unlockOk == 1
}

// UnlockWithContext Redis lock - release
//
// Use a Lua script to atomically verify the ownership token and delete the lock
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

	//Clean up in-memory data and stop automatic renewal
	rl.clearListenerAndStopAutoRenewal(key, ctrl)

	return unlockOk == 1
}

// Automatic renewal
func (rl *redisLockerImpl) autoRenewal(key, token string) {
	ctx, cancel := context.WithCancel(context.Background())
	ctrl := &cancelControl{ctx: ctx, cancel: cancel, token: token}

	states.mux.Lock()
	states.listeners[key] = ctrl
	states.mux.Unlock()
}

// Clean up the watcher and stop automatic renewal
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

// Central renewal scheduler
//
// 1.Use two-phase Mutex locking to reduce lock hold time
//
// 2.Use Redis pipelines to reduce RTT
func (rl *redisLockerImpl) startRenewalScheduler() {
	doRenewalRound := func() {
		if rl.client == nil {
			fmt.Println("[Go-Sail] <redisLock> renewal task not emit cause redis client is nil")
			return
		}

		states.mux.Lock()
		if len(states.listeners) == 0 {
			states.mux.Unlock()
			//Avoid holding locks while idle
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
			return //Handle leniently: skip cancel and delete in this round
		}

		states.mux.Lock()
		for _, command := range validCommands {
			expResult, expErr := command.cmd.Int64()
			if expResult != 1 || expErr != nil {
				if expErr != nil {
					fmt.Printf("[Go-Sail] <redisLock> key: %s renewal err: %v\n", command.key, expErr)
				}
				if current, ok := states.listeners[command.key]; ok && current == command.ctrl {
					current.cancel() //Clean up even if renewal fails
					delete(states.listeners, command.key)
				}
			}
		}
		//Clean up expired entries
		for index := range invalidKeys {
			if current, ok := states.listeners[invalidKeys[index].key]; ok && current == invalidKeys[index].ctrl {
				current.cancel() //Call again as a safeguard to ensure ctx.Done is triggered
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

// Precheck whether the locking task can run
//
// This is a local (in-memory) check
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

// Redis operation timeout control
func withRedisExecuteTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), redisExecuteTimeout)
}

var (
	hostname, _ = os.Hostname()   //Hostname
	ip, _       = IP().GetLocal() //Host IP
	processId   = os.Getpid()     //Process ID
)

// Lock holder information
func lockerValue() string {
	return fmt.Sprintf("lockedAt:%s@%s<%d>(%s:%s)",
		hostname, ip, processId, time.Now().Format("2006-01-02T15:04:05.000000Z"), rand.Text())
}
