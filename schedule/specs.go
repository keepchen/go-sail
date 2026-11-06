package schedule

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

// RunAt runs at a specified time
//
// Run tasks repeatedly at a specified time
//
// crontabExpr supports five- and six-field expressions; omitted seconds default to 0.
//
// Format: [second(0-59)] minute(0-59) hour(0-23) day(1-31) month(1-12) weekday(0-6, Sunday is 0).
// Months and weekdays also support English abbreviations; descriptors such as @daily and @every are supported.
// LastDayOfMonth is an additional month-end preset; other expressions containing L are not supported.
func (j *taskJob) RunAt(crontabExpr string) (cancel CancelFunc) {
	cronStartOnce.Do(func() {
		cronJob = cron.New(cron.WithParser(scheduleParser{}))
		cronJob.Start()
	})

	//Since AddFunc starts a goroutine internally, call the method here synchronously
	j.wrappedTaskFunc = func() {
		j.running = true

		defer func() {
			j.running = false
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

	if j.runOnceImmediately {
		go j.wrappedTaskFunc()
	}

	jobID, jobErr := cronJob.AddFunc(crontabExpr, j.wrappedTaskFunc)
	if jobErr != nil {
		fmt.Printf("[GO-SAIL] <Schedule> add job {%s} failed: %v\n", j.name, jobErr.Error())
	}

	cancel = func() {
		go func() {
			cronJob.Remove(jobID)
			taskSchedules.mux.Lock()
			delete(taskSchedules.pool, j.lockerKey)
			taskSchedules.mux.Unlock()
			fmt.Printf("[GO-SAIL] <Schedule> cancel job {%s} successfully\n", j.name)
		}()
	}

	return
}

// TenClockAtWeekday runs at 10:00 every weekday (Monday-Friday)
func (j *taskJob) TenClockAtWeekday() (cancel CancelFunc) {
	return j.RunAt(TenClockAtWeekday)
}

// TenClockAtWeekend runs at 10:00 every weekend day (Saturday and Sunday)
func (j *taskJob) TenClockAtWeekend() (cancel CancelFunc) {
	return j.RunAt(TenClockAtWeekend)
}

// FirstDayOfMonthly runs on the first day of every month
func (j *taskJob) FirstDayOfMonthly() (cancel CancelFunc) {
	return j.RunAt(FirstDayOfMonth)
}

// LastDayOfMonthly runs on the last day of every month
func (j *taskJob) LastDayOfMonthly() (cancel CancelFunc) {
	return j.RunAt(LastDayOfMonth)
}

// FirstDayOfWeek runs at 00:00 every Monday
func (j *taskJob) FirstDayOfWeek() (cancel CancelFunc) {
	return j.RunAt(FirstDayOfWeek)
}

// LastDayOfWeek runs at 00:00 every Sunday
func (j *taskJob) LastDayOfWeek() (cancel CancelFunc) {
	return j.RunAt(LastDayOfWeek)
}
