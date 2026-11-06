package schedule

import (
	"time"

	"github.com/keepchen/go-sail/v3/utils"
)

// setInterval sets the execution interval
//
// Note: interval must be at least 1 millisecond; otherwise it is set to 1 millisecond
func (j *taskJob) setInterval(interval time.Duration) *taskJob {
	if interval.Milliseconds() < 1 {
		interval = time.Millisecond
	}
	j.interval = interval

	return j
}

// Task execution function
func (j *taskJob) run() {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()

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

	LISTEN:
		for {
			select {
			case <-ticker.C:
				go j.wrappedTaskFunc()
			//Exit signal received; terminate the task
			case <-j.cancelTaskChan:
				if j.withoutOverlapping && j.lockedByMe {
					utils.RedisLocker(redisClients...).Unlock(j.lockerKey)
				}

				taskSchedules.mux.Lock()
				delete(taskSchedules.pool, j.lockerKey)
				taskSchedules.mux.Unlock()

				break LISTEN
			}
		}
	}()
}

// Every runs at the specified interval
//
// Note: interval must be at least 1 millisecond; otherwise it is set to 1 millisecond
func (j *taskJob) Every(interval time.Duration) (cancel CancelFunc) {
	j.setInterval(interval).run()

	cancel = j.cancelFunc

	return
}

// EverySecond runs every second
func (j *taskJob) EverySecond() (cancel CancelFunc) {
	j.setInterval(time.Second).run()

	cancel = j.cancelFunc

	return
}

// EveryFiveSeconds runs every 5 seconds
func (j *taskJob) EveryFiveSeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 5).run()

	cancel = j.cancelFunc

	return
}

// EveryTenSeconds runs every 10 seconds
func (j *taskJob) EveryTenSeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 10).run()

	cancel = j.cancelFunc

	return
}

// EveryFifteenSeconds runs every 15 seconds
func (j *taskJob) EveryFifteenSeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 15).run()

	cancel = j.cancelFunc

	return
}

// EveryTwentySeconds runs every 20 seconds
func (j *taskJob) EveryTwentySeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 20).run()

	cancel = j.cancelFunc

	return
}

// EveryThirtySeconds runs every 30 seconds
func (j *taskJob) EveryThirtySeconds() (cancel CancelFunc) {
	j.setInterval(time.Second * 30).run()

	cancel = j.cancelFunc

	return
}

// EveryMinute runs every minute
func (j *taskJob) EveryMinute() (cancel CancelFunc) {
	j.setInterval(time.Minute).run()

	cancel = j.cancelFunc

	return
}

// EveryFiveMinutes runs every 5 minutes
func (j *taskJob) EveryFiveMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 5).run()

	cancel = j.cancelFunc

	return
}

// EveryFifteenMinutes runs every 15 minutes
func (j *taskJob) EveryFifteenMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 15).run()

	cancel = j.cancelFunc

	return
}

// EveryTenMinutes runs every 10 minutes
func (j *taskJob) EveryTenMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 10).run()

	cancel = j.cancelFunc

	return
}

// EveryTwentyMinutes runs every 20 minutes
func (j *taskJob) EveryTwentyMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 20).run()

	cancel = j.cancelFunc

	return
}

// EveryThirtyMinutes runs every 30 minutes
func (j *taskJob) EveryThirtyMinutes() (cancel CancelFunc) {
	j.setInterval(time.Minute * 30).run()

	cancel = j.cancelFunc

	return
}

// Hourly runs every hour
func (j *taskJob) Hourly() (cancel CancelFunc) {
	j.setInterval(time.Hour).run()

	cancel = j.cancelFunc

	return
}

// EveryFiveHours runs every 5 hours
func (j *taskJob) EveryFiveHours() (cancel CancelFunc) {
	j.setInterval(time.Hour * 5).run()

	cancel = j.cancelFunc

	return
}

// EveryTenHours runs every 10 hours
func (j *taskJob) EveryTenHours() (cancel CancelFunc) {
	j.setInterval(time.Hour * 10).run()

	cancel = j.cancelFunc

	return
}

// EveryTwentyHours runs every 20 hours
func (j *taskJob) EveryTwentyHours() (cancel CancelFunc) {
	j.setInterval(time.Hour * 20).run()

	cancel = j.cancelFunc

	return
}

// Daily runs every day (every 24 hours)
func (j *taskJob) Daily() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24).run()

	cancel = j.cancelFunc

	return
}

// Weekly runs every week (every 7 days)
func (j *taskJob) Weekly() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24 * 7).run()

	cancel = j.cancelFunc

	return
}

// Monthly runs every month (every 30 days)
func (j *taskJob) Monthly() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24 * 30).run()

	cancel = j.cancelFunc

	return
}

// Yearly runs every year (every 365 days)
func (j *taskJob) Yearly() (cancel CancelFunc) {
	j.setInterval(time.Hour * 24 * 365).run()

	cancel = j.cancelFunc

	return
}
