package schedule

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

// RunAt 在某一时刻执行
//
// crontabExpr Linux crontab风格的表达式
//
// *    *    *    *    *
//
// -    -    -    -    -
//
// |    |    |    |    |
//
// |    |    |    |    +----- day of week (0 - 7) (Sunday=0 or 7) OR sun...sat
//
// |    |    |    +---------- month (1 - 12) OR jan,feb,mar,apr ...
//
// |    |    +--------------- day of month (1 - 31)
//
// |    +-------------------- hour (0 - 23)
//
// +------------------------- minute (0 - 59)
func (j *taskJob) RunAt(crontabExpr string) (cancel CancelFunc) {
	cronStartOnce.Do(func() {
		cronJob = cron.New()
		cronJob.Start()
	})

	//因为AddFunc内部是协程启动，因此这里的方法使用同步方式调用
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

	jobID, jobErr := cronJob.AddFunc(crontabExpr, wrappedTaskFunc)
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

// TenClockAtWeekday 每个工作日（周一~周五）上午10点
func (j *taskJob) TenClockAtWeekday() (cancel CancelFunc) {
	return j.RunAt(TenClockAtWeekday)
}

// TenClockAtWeekend 每个周末（周六和周日）上午10点
func (j *taskJob) TenClockAtWeekend() (cancel CancelFunc) {
	return j.RunAt(TenClockAtWeekend)
}

// FirstDayOfMonthly 每月1号
func (j *taskJob) FirstDayOfMonthly() (cancel CancelFunc) {
	return j.RunAt(FirstDayOfMonth)
}

// LastDayOfMonthly 每月最后一天
func (j *taskJob) LastDayOfMonthly() (cancel CancelFunc) {
	return j.RunAt(LastDayOfMonth)
}

// FirstDayOfWeek 每周1的00:00
func (j *taskJob) FirstDayOfWeek() (cancel CancelFunc) {
	return j.RunAt(FirstDayOfWeek)
}

// LastDayOfWeek 每周天的00:00
func (j *taskJob) LastDayOfWeek() (cancel CancelFunc) {
	return j.RunAt(LastDayOfWeek)
}
