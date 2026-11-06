package schedule

import (
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// RunAfter runs after a specified duration
//
// # Note
//
// This is a one-time task; it runs once and does not repeat
func (j *taskJob) RunAfter(delay time.Duration) (cancel CancelFunc) {
	timer := time.After(delay)
	cancel = j.cancelFunc

	j.mux.Lock()
	j.wrappedTaskFunc = func() {
		j.runningMux.Lock()
		j.running = true
		j.runningMux.Unlock()

		defer func() {
			j.runningMux.Lock()
			j.running = false
			j.runningMux.Unlock()
		}()

		if !j.withoutOverlapping {
			j.task()
			return
		}
		if utils.RedisLocker(redisClients...).TryLock(j.lockerKey) {
			defer func() {
				utils.RedisLocker(redisClients...).Unlock(j.lockerKey)
				j.lockedByMe = false
			}()
			j.lockedByMe = true
			j.task()
		}
	}
	j.mux.Unlock()

	if j.runOnceImmediately {
		go j.wrappedTaskFunc()
	}

	go func() {
	LOOP:
		for {
			select {
			case <-timer:
				go j.wrappedTaskFunc()
				break LOOP
			case <-j.cancelTaskChan:
				break LOOP
			}
		}
	}()

	return cancel
}
