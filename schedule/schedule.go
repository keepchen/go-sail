package schedule

import (
	"fmt"
	"sync"
	"time"

	redisLib "github.com/go-redis/redis/v8"

	"github.com/keepchen/go-sail/v3/utils"
	"github.com/robfig/cron/v3"
)

// CancelFunc cancellation function
//
// Tasks that have not started or are not running are canceled immediately.
// Running tasks are allowed to finish and will not run again.
//
// Calling this method removes the task from the task list.
type CancelFunc func()

// Scheduler scheduler
type Scheduler interface {
	// WithoutOverlapping prevents task overlap by disallowing concurrent execution
	//
	// Only one execution of a task may be running at a time
	//
	// # Note: this method uses a Redis lock to ensure 'global' uniqueness,
	//
	// so initialize the Redis connection first with redis.InitRedis or
	//
	// redis.InitRedisCluster
	WithoutOverlapping() Scheduler
	// RunOnceImmediately runs once immediately
	//
	// Run once immediately after successful registration, before the scheduled trigger
	//
	// # Note:
	//
	// Still respects WithoutOverlapping (if enabled)
	RunOnceImmediately() Scheduler
	// Every runs at the specified interval
	//
	// Note: interval must be at least 1 millisecond; otherwise it is set to 1 millisecond
	Every(interval time.Duration) (cancel CancelFunc)
	// EverySecond runs every second
	EverySecond() (cancel CancelFunc)
	// EveryFiveSeconds runs every 5 seconds
	EveryFiveSeconds() (cancel CancelFunc)
	// EveryTenSeconds runs every 10 seconds
	EveryTenSeconds() (cancel CancelFunc)
	// EveryFifteenSeconds runs every 15 seconds
	EveryFifteenSeconds() (cancel CancelFunc)
	// EveryTwentySeconds runs every 20 seconds
	EveryTwentySeconds() (cancel CancelFunc)
	// EveryThirtySeconds runs every 30 seconds
	EveryThirtySeconds() (cancel CancelFunc)
	// EveryMinute runs every minute
	EveryMinute() (cancel CancelFunc)
	// EveryFiveMinutes runs every 5 minutes
	EveryFiveMinutes() (cancel CancelFunc)
	// EveryTenMinutes runs every 10 minutes
	EveryTenMinutes() (cancel CancelFunc)
	// EveryFifteenMinutes runs every 15 minutes
	EveryFifteenMinutes() (cancel CancelFunc)
	// EveryTwentyMinutes runs every 20 minutes
	EveryTwentyMinutes() (cancel CancelFunc)
	// EveryThirtyMinutes runs every 30 minutes
	EveryThirtyMinutes() (cancel CancelFunc)
	// Hourly runs every hour
	Hourly() (cancel CancelFunc)
	// EveryFiveHours runs every 5 hours
	EveryFiveHours() (cancel CancelFunc)
	// EveryTenHours runs every 10 hours
	EveryTenHours() (cancel CancelFunc)
	// EveryTwentyHours runs every 20 hours
	EveryTwentyHours() (cancel CancelFunc)
	// Daily runs every day (every 24 hours)
	Daily() (cancel CancelFunc)
	// Weekly runs every week (every 7 days)
	Weekly() (cancel CancelFunc)
	// Monthly runs every month (every 30 days)
	Monthly() (cancel CancelFunc)
	// Yearly runs every year (every 365 days)
	Yearly() (cancel CancelFunc)
	// RunAfter runs after a specified duration
	//
	// # Note
	//
	// This is a one-time task and does not repeat
	RunAfter(delay time.Duration) (cancel CancelFunc)
	// RunAt runs at a specified time
	//
	//Run tasks repeatedly at a specified time
	//
	// crontabExpr supports five- and six-field expressions; omitted seconds default to 0.
	//
	// Format: [second(0-59)] minute(0-59) hour(0-23) day(1-31) month(1-12) weekday(0-6, Sunday is 0).
	// Months and weekdays also support English abbreviations; descriptors such as @daily and @every are supported.
	// LastDayOfMonth is an additional month-end preset; other expressions containing L are not supported.
	RunAt(crontabExpr string) (cancel CancelFunc)
	// TenClockAtWeekday runs at 10:00 every weekday (Monday-Friday)
	TenClockAtWeekday() (cancel CancelFunc)
	// TenClockAtWeekend runs at 10:00 every weekend day (Saturday and Sunday)
	TenClockAtWeekend() (cancel CancelFunc)
	// FirstDayOfMonthly runs on the first day of every month
	FirstDayOfMonthly() (cancel CancelFunc)
	// LastDayOfMonthly runs on the last day of every month
	LastDayOfMonthly() (cancel CancelFunc)
	// FirstDayOfWeek runs at 00:00 every Monday
	FirstDayOfWeek() (cancel CancelFunc)
	// LastDayOfWeek runs at 00:00 every Sunday
	LastDayOfWeek() (cancel CancelFunc)
}

// taskJob task
type taskJob struct {
	name               string
	task               func()
	wrappedTaskFunc    func()
	interval           time.Duration
	lockerKey          string
	lockedByMe         bool
	running            bool
	runningMux         sync.RWMutex
	withoutOverlapping bool
	cancelFunc         CancelFunc
	cancelTaskChan     chan struct{}
	cancelOnce         sync.Once
	mux                sync.RWMutex
	runOnceImmediately bool
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
	return fmt.Sprintf("go-sail:task-schedule-locker:%s", utils.Base64().Encode([]byte(name)))
}

var (
	setRedisOnce = &sync.Once{}
	redisClients = make([]redisLib.UniversalClient, 0)
)

// SetRedisClientOnce sets the Redis connection client
func SetRedisClientOnce(client redisLib.UniversalClient) {
	setRedisOnce.Do(func() {
		redisClients = []redisLib.UniversalClient{client}
	})
}

// NewJob instantiates a task
//
// name unique task name identifier
//
// task task handler
//
// # Note:
//
// Panics if name is duplicated
func NewJob(name string, task func()) Scheduler {
	return Job(name, task)
}

// Job instantiates a task
//
// name unique task name identifier
//
// task task handler
//
// # Note:
//
// Panics if name is duplicated
func Job(name string, task func()) Scheduler {
	job := &taskJob{
		name:           name,
		lockerKey:      generateJobNameKey(name),
		task:           task,
		cancelTaskChan: make(chan struct{}),
		cancelOnce:     sync.Once{},
	}

	job.cancelFunc = func() {
		job.cancelOnce.Do(func() {
			go func() {
				job.cancelTaskChan <- struct{}{}
				close(job.cancelTaskChan)
				taskSchedules.mux.Lock()
				delete(taskSchedules.pool, job.lockerKey)
				taskSchedules.mux.Unlock()
				fmt.Printf("[GO-SAIL] <Schedule> cancel job {%s} successfully\n", job.name)
			}()
		})
	}

	taskSchedules.mux.Lock()
	if _, ok := taskSchedules.pool[job.lockerKey]; !ok {
		taskSchedules.pool[job.lockerKey] = job
	} else {
		taskSchedules.mux.Unlock()
		panic(fmt.Errorf("[GO-SAIL] Duplicate schedule, task name: %s", job.name))
	}
	taskSchedules.mux.Unlock()

	return job
}

// WithoutOverlapping prevents task overlap by disallowing concurrent execution
//
// Only one execution of a task may be running at a time
//
// # Note: this method uses a Redis lock to ensure 'global' uniqueness,
//
// so initialize the Redis connection first with redis.InitRedis or
//
// redis.InitRedisCluster
func (j *taskJob) WithoutOverlapping() Scheduler {
	j.withoutOverlapping = true

	return j
}

// RunOnceImmediately runs once immediately
//
// Run once immediately after successful registration, before the scheduled trigger
//
// # Note:
//
// Still respects WithoutOverlapping (if enabled)
func (j *taskJob) RunOnceImmediately() Scheduler {
	j.runOnceImmediately = true

	return j
}

// JobIsRunning checks whether a task is running
func JobIsRunning(jobName string) bool {
	var (
		running = false
		name    = generateJobNameKey(jobName)
	)
	taskSchedules.mux.RLock()
	if job, ok := taskSchedules.pool[name]; ok {
		job.runningMux.RLock()
		running = job.running
		job.runningMux.RUnlock()
	}
	taskSchedules.mux.RUnlock()

	return running
}

// Call starts a task manually
//
// jobName task name
//
// mandatory if true, execute directly without checking for overlap
//
// # Note
//
// The internal function is called synchronously
func Call(jobName string, mandatory bool) {
	var (
		job       *taskJob
		lockerKey = generateJobNameKey(jobName)
	)
	taskSchedules.mux.RLock()
	if jb, ok := taskSchedules.pool[lockerKey]; ok {
		job = jb
	}
	taskSchedules.mux.RUnlock()
	if job == nil {
		fmt.Printf("[GO-SAIL] <Schedule> call job {%s} failed,cause job not registered.\n", jobName)
		return
	}
	if mandatory {
		job.task()
	} else {
		ticker := time.NewTicker(time.Millisecond * 500)
		retryTimes := 0
		defer ticker.Stop()
	RetryLoop:
		for {
			select {
			case <-ticker.C:
				job.mux.RLock()
				wrappedTaskFunc := job.wrappedTaskFunc
				job.mux.RUnlock()
				if wrappedTaskFunc != nil {
					wrappedTaskFunc()
					break RetryLoop
				}
				retryTimes++

				if retryTimes > 2 {
					break RetryLoop
				}
			}
		}
	}
}

// MustCall starts a task manually
//
// jobName task name
//
// mandatory if true, execute directly without checking for overlap
//
// # Note
//
// 1.Panics if jobName is absent from the task list (e.g. not registered or canceled)
//
// 2.The internal function is called synchronously
func MustCall(jobName string, mandatory bool) {
	var (
		job       *taskJob
		lockerKey = generateJobNameKey(jobName)
	)
	taskSchedules.mux.RLock()
	if jb, ok := taskSchedules.pool[lockerKey]; ok {
		job = jb
	}
	taskSchedules.mux.RUnlock()
	if job == nil {
		panic(fmt.Errorf("job name: %s not registered", jobName))
	}
	if mandatory {
		job.task()
	} else {
		ticker := time.NewTicker(time.Millisecond * 500)
		retryTimes := 0
		defer ticker.Stop()
	RetryLoop:
		for {
			select {
			case <-ticker.C:
				job.mux.RLock()
				wrappedTaskFunc := job.wrappedTaskFunc
				job.mux.RUnlock()
				if wrappedTaskFunc != nil {
					wrappedTaskFunc()
					break RetryLoop
				}
				retryTimes++

				if retryTimes > 2 {
					break RetryLoop
				}
			}
		}
	}
}
