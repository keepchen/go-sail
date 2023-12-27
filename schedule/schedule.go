package schedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

// TaskJob 任务
type TaskJob struct {
	name               string
	task               func()
	interval           time.Duration
	lockerKey          string
	cancelTaskChan     chan struct{}
	withoutOverlapping bool
}

var (
	cancelFunc func()
	cronJob    *cron.Cron
)

// Job 实例化任务
//
// @param name 任务名称唯一标识
//
// @param task 任务处理函数
//
// Note: 如果需要保证任务同一时刻只有一个运行态，task内部不要使用协程运行主逻辑。
func Job(name string, task func()) *TaskJob {
	job := &TaskJob{
		name:           name,
		lockerKey:      fmt.Sprintf("go-sail:task-schedule-locker:%s", utils.MD5Encrypt(name)),
		task:           task,
		cancelTaskChan: make(chan struct{}),
	}

	cancelFunc = func() {
		go func() {
			job.cancelTaskChan <- struct{}{}
			close(job.cancelTaskChan)
			fmt.Printf("[go-sail] task schedule cancel job {%s} successfully\n", job.name)
		}()
	}

	return job
}

// WithoutOverlapping 禁止并发执行
//
// 同一时刻仅允许一个任务线程执行
//
// Note: 该方法使用redis锁来保证唯一性
func (j *TaskJob) WithoutOverlapping() *TaskJob {
	j.withoutOverlapping = true

	return j
}

// Every 每隔多久执行一次
func (j *TaskJob) Every(interval time.Duration) (cancel func()) {
	j.interval = interval
	j.run()

	cancel = cancelFunc

	return cancel
}

// EverySecond 每秒执行一次
func (j *TaskJob) EverySecond() (cancel func()) {
	j.interval = time.Second
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryFiveSeconds 每5秒执行一次
func (j *TaskJob) EveryFiveSeconds() (cancel func()) {
	j.interval = time.Second * 5
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryTenSeconds 每10秒执行一次
func (j *TaskJob) EveryTenSeconds() (cancel func()) {
	j.interval = time.Second * 10
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryThirtySeconds 每30秒执行一次
func (j *TaskJob) EveryThirtySeconds() (cancel func()) {
	j.interval = time.Second * 30
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryMinute 每分钟执行一次
func (j *TaskJob) EveryMinute() (cancel func()) {
	j.interval = time.Minute
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryFiveMinutes 每5分钟执行一次
func (j *TaskJob) EveryFiveMinutes() (cancel func()) {
	j.interval = time.Minute * 5
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryTenMinutes 每10分钟执行一次
func (j *TaskJob) EveryTenMinutes() (cancel func()) {
	j.interval = time.Minute * 10
	j.run()

	cancel = cancelFunc

	return cancel
}

// EveryThirtyMinutes 每30分钟执行一次
func (j *TaskJob) EveryThirtyMinutes() (cancel func()) {
	j.interval = time.Minute * 30
	j.run()

	cancel = cancelFunc

	return cancel
}

// Hourly 每1小时执行一次
func (j *TaskJob) Hourly() (cancel func()) {
	j.interval = time.Hour
	j.run()

	cancel = cancelFunc

	return cancel
}

// Daily 每天执行一次
func (j *TaskJob) Daily() (cancel func()) {
	j.interval = time.Hour * 24
	j.run()

	cancel = cancelFunc

	return cancel
}

// Weekly 每周执行一次（每7天）
func (j *TaskJob) Weekly() (cancel func()) {
	j.interval = time.Hour * 24 * 7
	j.run()

	cancel = cancelFunc

	return cancel
}

// Monthly 每月执行一次（每30天）
func (j *TaskJob) Monthly() (cancel func()) {
	j.interval = time.Hour * 24 * 30
	j.run()

	cancel = cancelFunc

	return cancel
}

// Yearly 每年执行一次（每365天）
func (j *TaskJob) Yearly() (cancel func()) {
	j.interval = time.Hour * 24 * 365
	j.run()

	cancel = cancelFunc

	return cancel
}

// 任务执行函数
func (j *TaskJob) run() {
	go func() {
		ticker := time.NewTicker(j.interval)
		for {
			select {
			case <-ticker.C:
				if j.withoutOverlapping {
					if utils.RedisLock(j.lockerKey) {
						go func() {
							defer utils.RedisUnlock(j.lockerKey)
							j.task()
						}()
					}
				} else {
					go j.task()
				}
			//shutdown gracefully.
			case <-j.cancelTaskChan:
				utils.RedisUnlock(j.lockerKey)
				ticker.Stop()
				break
			}
		}
	}()
}

// RunAt 在某一时刻执行
//
// @param crontabExpr Linux crontab风格的表达式
//
// *    *    *    *    *
//
// -    -    -    -    -
//
// |    |    |    |    |
//
// |    |    |    |    +----- day of week (0 - 7) (Sunday=0 or 7) OR sun,mon,tue,wed,thu,fri,sat
//
// |    |    |    +---------- month (1 - 12) OR jan,feb,mar,apr ...
//
// |    |    +--------------- day of month (1 - 31)
//
// |    +-------------------- hour (0 - 23)
//
// +------------------------- minute (0 - 59)
func (j *TaskJob) RunAt(crontabExpr string) (cancel func()) {
	(&sync.Once{}).Do(func() {
		cronJob = cron.New()
		cronJob.Start()
	})

	wrappedTaskFunc := func() {
		if j.withoutOverlapping {
			if utils.RedisLock(j.lockerKey) {
				go func() {
					defer utils.RedisUnlock(j.lockerKey)
					j.task()
				}()
			}
		} else {
			j.task()
		}
	}

	jobID, jobErr := cronJob.AddFunc(crontabExpr, wrappedTaskFunc)
	if jobErr != nil {
		fmt.Printf("[go-sail] task schedule add job failed: %v\n", jobErr.Error())
	}

	cancel = func() {
		go func() {
			cronJob.Remove(jobID)
			fmt.Printf("[go-sail] task schedule cancel job {%s} successfully\n", j.name)
		}()
	}

	return
}
