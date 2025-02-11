package schedule

import (
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// setInterval 设置执行间隔
//
// Note: interval至少需要大于等于1毫秒，否则将被设置为1毫秒
func (j *taskJob) setInterval(interval time.Duration) *taskJob {
	if interval.Milliseconds() < 1 {
		interval = time.Millisecond
	}
	j.interval = interval

	return j
}

// 任务执行函数
func (j *taskJob) run() {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()
		j.wrappedTaskFunc = func() {
			j.running = true

			defer func() {
				j.running = false
			}()

			if !j.withoutOverlapping {
				j.task()
				return
			}
			if utils.RedisLocker().TryLock(j.lockerKey) {
				defer func() {
					utils.RedisLocker().Unlock(j.lockerKey)
					j.lockedByMe = false
				}()
				j.lockedByMe = true
				j.task()
			}
		}
	LISTEN:
		for {
			select {
			case <-ticker.C:
				go j.wrappedTaskFunc()
			//收到退出信号，终止任务
			case <-j.cancelTaskChan:
				if j.withoutOverlapping && j.lockedByMe {
					utils.RedisLocker().Unlock(j.lockerKey)
				}

				taskSchedules.mux.Lock()
				delete(taskSchedules.pool, j.lockerKey)
				taskSchedules.mux.Unlock()

				break LISTEN
			}
		}
	}()
}

// Every 每隔多久执行一次
//
// Note: interval至少需要大于等于1毫秒，否则将被设置为1毫秒
func (j *taskJob) Every(interval time.Duration) (cancel CancelFunc) {
	j.setInterval(interval).run()

	cancel = j.cancelFunc

	return
}

// EverySecond 每秒执行一次
func (j *taskJob) EverySecond() (cancel CancelFunc) {
	j.setInterval(time.Second).run()

	cancel = j.cancelFunc

	return
}

// EveryFiveSeconds 每5秒执行一次
func (j *taskJob) EveryFiveSeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 5).run()

	cancel = j.cancelFunc

	return
}

// EveryTenSeconds 每10秒执行一次
func (j *taskJob) EveryTenSeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 10).run()

	cancel = j.cancelFunc

	return
}

// EveryFifteenSeconds 每15秒执行一次
func (j *taskJob) EveryFifteenSeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 15).run()

	cancel = j.cancelFunc

	return
}

// EveryTwentySeconds 每20秒执行一次
func (j *taskJob) EveryTwentySeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 20).run()

	cancel = j.cancelFunc

	return
}

// EveryThirtySeconds 每30秒执行一次
func (j *taskJob) EveryThirtySeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 30).run()

	cancel = j.cancelFunc

	return
}

// EveryMinute 每分钟执行一次
func (j *taskJob) EveryMinute() (cancel CancelFunc) {
	j.setInterval(time.Minute).run()

	cancel = j.cancelFunc

	return
}

// EveryFiveMinutes 每5分钟执行一次
func (j *taskJob) EveryFiveMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 5).run()

	cancel = j.cancelFunc

	return
}

// EveryFifteenMinutes 每15分钟执行一次
func (j *taskJob) EveryFifteenMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 15).run()

	cancel = j.cancelFunc

	return
}

// EveryTenMinutes 每10分钟执行一次
func (j *taskJob) EveryTenMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 10).run()

	cancel = j.cancelFunc

	return
}

// EveryTwentyMinutes 每20分钟执行一次
func (j *taskJob) EveryTwentyMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 20).run()

	cancel = j.cancelFunc

	return
}

// EveryThirtyMinutes 每30分钟执行一次
func (j *taskJob) EveryThirtyMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 30).run()

	cancel = j.cancelFunc

	return
}

// Hourly 每1小时执行一次
func (j *taskJob) Hourly() (cancel CancelFunc) {
	j.setInterval(time.Hour).run()

	cancel = j.cancelFunc

	return
}

// EveryFiveHours 每5小时执行一次
func (j *taskJob) EveryFiveHours() (cancel CancelFunc) {
	j.setInterval(time.Hour * 5).run()

	cancel = j.cancelFunc

	return
}

// EveryTenHours 每10小时执行一次
func (j *taskJob) EveryTenHours() (cancel CancelFunc) {
	j.setInterval(time.Hour * 10).run()

	cancel = j.cancelFunc

	return
}

// EveryTwentyHours 每20小时执行一次
func (j *taskJob) EveryTwentyHours() (cancel CancelFunc) {
	j.setInterval(time.Hour * 20).run()

	cancel = j.cancelFunc

	return
}

// Daily 每天执行一次（每24小时）
func (j *taskJob) Daily() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24).run()

	cancel = j.cancelFunc

	return
}

// Weekly 每周执行一次（每7天）
func (j *taskJob) Weekly() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24 * 7).run()

	cancel = j.cancelFunc

	return
}

// Monthly 每月执行一次（每30天）
func (j *taskJob) Monthly() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24 * 30).run()

	cancel = j.cancelFunc

	return
}

// Yearly 每年执行一次（每365天）
func (j *taskJob) Yearly() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24 * 365).run()

	cancel = j.cancelFunc

	return
}
