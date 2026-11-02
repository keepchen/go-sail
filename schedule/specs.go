package schedule

import (
	"fmt"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

// RunAt 在某一时刻执行
//
// 重复地在某个时间点执行任务
//
// crontabExpr 支持五字段和六字段表达式，省略秒字段时默认为0。
//
// 格式：[秒(0-59)] 分(0-59) 时(0-23) 日(1-31) 月(1-12) 星期(0-6，周日为0)。
// 月份和星期也支持英文缩写；支持 @daily、@every 等描述符。
// LastDayOfMonth 是额外支持的月末预设，不支持其他含 L 的表达式。
func (j *taskJob) RunAt(crontabExpr string) (cancel CancelFunc) {
	cronStartOnce.Do(func() {
		cronJob = cron.New(cron.WithParser(scheduleParser{}))
		cronJob.Start()
	})

	//因为AddFunc内部是协程启动，因此这里的方法使用同步方式调用
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
