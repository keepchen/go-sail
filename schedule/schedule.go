package schedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

type CancelFunc func()

// Scheduler 调度器
type Scheduler interface {
	// WithoutOverlapping 禁止并发执行
	//
	// 一个任务仅允许存在一个运行态
	//
	// Note: 该方法使用redis锁来保证唯一性，
	//
	// 因此请确保先使用 redis.InitRedis 或
	//
	// redis.InitRedisCluster 实例化redis连接
	WithoutOverlapping() Scheduler
	//任务执行函数
	run()
	// Every 每隔多久执行一次
	//
	// Note: interval至少需要大于等于1毫秒，否则将被设置为1毫秒
	Every(interval time.Duration) (cancel CancelFunc)
	// EverySecond 每秒执行一次
	EverySecond() (cancel CancelFunc)
	// EveryFiveSeconds 每5秒执行一次
	EveryFiveSeconds() (cancel CancelFunc)
	// EveryTenSeconds 每10秒执行一次
	EveryTenSeconds() (cancel CancelFunc)
	// EveryTwentySeconds 每20秒执行一次
	EveryTwentySeconds() (cancel CancelFunc)
	// EveryThirtySeconds 每30秒执行一次
	EveryThirtySeconds() (cancel CancelFunc)
	// EveryMinute 每分钟执行一次
	EveryMinute() (cancel CancelFunc)
	// EveryFiveMinutes 每5分钟执行一次
	EveryFiveMinutes() (cancel CancelFunc)
	// EveryTenMinutes 每10分钟执行一次
	EveryTenMinutes() (cancel CancelFunc)
	// EveryTwentyMinutes 每20分钟执行一次
	EveryTwentyMinutes() (cancel CancelFunc)
	// EveryThirtyMinutes 每30分钟执行一次
	EveryThirtyMinutes() (cancel CancelFunc)
	// Hourly 每1小时执行一次
	Hourly() (cancel CancelFunc)
	// EveryFiveHours 每5小时执行一次
	EveryFiveHours() (cancel CancelFunc)
	// EveryTenHours 每10小时执行一次
	EveryTenHours() (cancel CancelFunc)
	// EveryTwentyHours 每20小时执行一次
	EveryTwentyHours() (cancel CancelFunc)
	// Daily 每天执行一次
	Daily() (cancel CancelFunc)
	// Weekly 每周执行一次（每7天）
	Weekly() (cancel CancelFunc)
	// Monthly 每月执行一次（每30天）
	Monthly() (cancel CancelFunc)
	// Yearly 每年执行一次（每365天）
	Yearly() (cancel CancelFunc)
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
	// |    |    |    |    +----- day of week (0 - 7) (Sunday=0 or 7) OR sun...sat
	//
	// |    |    |    +---------- month (1 - 12) OR jan,feb,mar,apr ...
	//
	// |    |    +--------------- day of month (1 - 31)
	//
	// |    +-------------------- hour (0 - 23)
	//
	// +------------------------- minute (0 - 59)
	RunAt(crontabExpr string) (cancel CancelFunc)
	// TenClockAtWeekday 每个工作日（周一~周五）上午10点
	TenClockAtWeekday() (cancel CancelFunc)
	// TenClockAtWeekend 每个周末（周六和周日）上午10点
	TenClockAtWeekend() (cancel CancelFunc)
	// FirstDayOfMonthly 每月1号
	FirstDayOfMonthly() (cancel CancelFunc)
	// LastDayOfMonthly 每月最后一天
	LastDayOfMonthly() (cancel CancelFunc)
}

// taskJob 任务
type taskJob struct {
	name               string
	task               func()
	interval           time.Duration
	lockerKey          string
	lockedByMe         bool
	running            bool
	withoutOverlapping bool
	cancelFunc         CancelFunc
	cancelTaskChan     chan struct{}
}

var _ Scheduler = &taskJob{}

type taskJobPool struct {
	mux  *sync.RWMutex
	pool map[string]*taskJob
}

var taskSchedules *taskJobPool

var cronJob *cron.Cron

func init() {
	(&sync.Once{}).Do(func() {
		taskSchedules = &taskJobPool{
			mux:  &sync.RWMutex{},
			pool: make(map[string]*taskJob),
		}
	})
}

func generateJobNameKey(name string) string {
	return fmt.Sprintf("go-sail:task-schedule-locker:%s", utils.MD5Encrypt(name))
}

// Job 实例化任务
//
// @param name 任务名称唯一标识
//
// @param task 任务处理函数
func Job(name string, task func()) Scheduler {
	job := &taskJob{
		name:           name,
		lockerKey:      generateJobNameKey(name),
		task:           task,
		cancelTaskChan: make(chan struct{}),
	}

	job.cancelFunc = func() {
		go func() {
			job.cancelTaskChan <- struct{}{}
			close(job.cancelTaskChan)
			fmt.Printf("[GO-SAIL] <Schedule> cancel job {%s} successfully\n", job.name)
		}()
	}

	taskSchedules.mux.Lock()
	taskSchedules.pool[job.lockerKey] = job
	taskSchedules.mux.Unlock()

	return job
}

// WithoutOverlapping 禁止并发执行
//
// 一个任务仅允许存在一个运行态
//
// Note: 该方法使用redis锁来保证唯一性，
//
// 因此请确保先使用 redis.InitRedis 或
//
// redis.InitRedisCluster 实例化redis连接
func (j *taskJob) WithoutOverlapping() Scheduler {
	j.withoutOverlapping = true

	return j
}

// 任务执行函数
func (j *taskJob) run() {
	go func() {
		ticker := time.NewTicker(j.interval)
		defer ticker.Stop()
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
	LISTEN:
		for {
			select {
			case <-ticker.C:
				go wrappedTaskFunc()
			//收到退出信号，终止任务
			case <-j.cancelTaskChan:
				if j.withoutOverlapping && j.lockedByMe {
					utils.RedisUnlock(j.lockerKey)
				}

				taskSchedules.mux.Lock()
				delete(taskSchedules.pool, j.lockerKey)
				taskSchedules.mux.Unlock()

				break LISTEN
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
	(&sync.Once{}).Do(func() {
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

// JobIsRunning 查看任务是否正在执行
func JobIsRunning(jobName string) bool {
	var (
		running = false
		name    = generateJobNameKey(jobName)
	)
	taskSchedules.mux.RLock()
	if job, ok := taskSchedules.pool[name]; ok {
		running = job.running
	}
	taskSchedules.mux.RUnlock()

	return running
}
