package schedule

import (
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// RunAfter 在一定时间后执行
//
// # Note
//
// 这是一个一次性任务，只会执行一次，不会重复执行
func (j *taskJob) RunAfter(delay time.Duration) (cancel CancelFunc) {
	timer := time.After(delay)
	cancel = j.cancelFunc

	wrappedTaskFunc := func() {
		j.running = true

		defer func() {
			j.running = false
		}()

		if !j.withoutOverlapping {
			j.task()
			return
		}
		if utils.RedisTryLock(j.lockerKey) {
			defer func() {
				utils.RedisUnlock(j.lockerKey)
				j.lockedByMe = false
			}()
			j.lockedByMe = true
			j.task()
		}
	}

	go func() {
	LOOP:
		for {
			select {
			case <-timer:
				go wrappedTaskFunc()
				break LOOP
			case <-j.cancelTaskChan:
				break LOOP
			}
		}
	}()

	return cancel
}
