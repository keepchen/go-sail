package schedule

import (
	"fmt"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

// CancelFunc 取消函数
//
// 尚未启动或未在运行中的任务将被直接取消。
// 正在运行中的任务将等待其运行结束，之后便不再启动。
type CancelFunc func()

// Scheduler 调度器
type Scheduler interface {
	// WithoutOverlapping 禁止并发执行
	//
	// 一个任务仅允许存在一个运行态
	//
	// Note: 该方法使用redis锁来保证'全局'唯一性，
	//
	// 因此请确保先使用 redis.InitRedis 或
	//
	// redis.InitRedisCluster 实例化redis连接
	WithoutOverlapping() Scheduler
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
	// EveryFifteenSeconds 每15秒执行一次
	EveryFifteenSeconds() (cancel CancelFunc)
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
	// EveryFifteenMinutes 每15分钟执行一次
	EveryFifteenMinutes() (cancel CancelFunc)
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
	// Daily 每天执行一次（每24小时）
	Daily() (cancel CancelFunc)
	// Weekly 每周执行一次（每7天）
	Weekly() (cancel CancelFunc)
	// Monthly 每月执行一次（每30天）
	Monthly() (cancel CancelFunc)
	// Yearly 每年执行一次（每365天）
	Yearly() (cancel CancelFunc)
	// RunAfter 在一定时间后执行
	//
	// # Note
	//
	// 这是一个一次性任务，不会重复执行
	RunAfter(delay time.Duration) (cancel CancelFunc)
	// RunAt 在某一时刻执行
	//
	//重复地在某个时间点执行任务
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
	RunAt(crontabExpr string) (cancel CancelFunc)
	// TenClockAtWeekday 每个工作日（周一~周五）上午10点
	TenClockAtWeekday() (cancel CancelFunc)
	// TenClockAtWeekend 每个周末（周六和周日）上午10点
	TenClockAtWeekend() (cancel CancelFunc)
	// FirstDayOfMonthly 每月1号
	FirstDayOfMonthly() (cancel CancelFunc)
	// LastDayOfMonthly 每月最后一天
	LastDayOfMonthly() (cancel CancelFunc)
	// FirstDayOfWeek 每周1的00:00
	FirstDayOfWeek() (cancel CancelFunc)
	// LastDayOfWeek 每周天的00:00
	LastDayOfWeek() (cancel CancelFunc)
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

var (
	cronJob       *cron.Cron
	cronStartOnce sync.Once
	initOnce      sync.Once
)

func init() {
	initOnce.Do(func() {
		taskSchedules = &taskJobPool{
			mux:  &sync.RWMutex{},
			pool: make(map[string]*taskJob),
		}
	})
}

func generateJobNameKey(name string) string {
	return fmt.Sprintf("go-sail:task-schedule-locker:%s", utils.MD5Encode(name))
}

// NewJob 实例化任务
//
// name 任务名称唯一标识
//
// task 任务处理函数
func NewJob(name string, task func()) Scheduler {
	return Job(name, task)
}

// Job 实例化任务
//
// name 任务名称唯一标识
//
// task 任务处理函数
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
// Note: 该方法使用redis锁来保证'全局'唯一性，
//
// 因此请确保先使用 redis.InitRedis 或
//
// redis.InitRedisCluster 实例化redis连接
func (j *taskJob) WithoutOverlapping() Scheduler {
	j.withoutOverlapping = true

	return j
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
