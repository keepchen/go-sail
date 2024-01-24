package schedule

import "time"

// Every 每隔多久执行一次
//
// Note: interval至少需要大于等于1毫秒，否则将被设置为1毫秒
func (j *TaskJob) Every(interval time.Duration) (cancel CancelFunc) {
	if interval.Milliseconds() < 1 {
		interval = time.Millisecond
	}
	j.interval = interval
	j.run()

	cancel = j.cancelFunc

	return
}

// EverySecond 每秒执行一次
func (j *TaskJob) EverySecond() (cancel CancelFunc) {
	j.interval = time.Second
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryFiveSeconds 每5秒执行一次
func (j *TaskJob) EveryFiveSeconds() (cancel CancelFunc) {
	j.interval = time.Second * 5
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryTenSeconds 每10秒执行一次
func (j *TaskJob) EveryTenSeconds() (cancel CancelFunc) {
	j.interval = time.Second * 10
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryTwentySeconds 每20秒执行一次
func (j *TaskJob) EveryTwentySeconds() (cancel CancelFunc) {
	j.interval = time.Second * 20
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryThirtySeconds 每30秒执行一次
func (j *TaskJob) EveryThirtySeconds() (cancel CancelFunc) {
	j.interval = time.Second * 30
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryMinute 每分钟执行一次
func (j *TaskJob) EveryMinute() (cancel CancelFunc) {
	j.interval = time.Minute
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryFiveMinutes 每5分钟执行一次
func (j *TaskJob) EveryFiveMinutes() (cancel CancelFunc) {
	j.interval = time.Minute * 5
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryTenMinutes 每10分钟执行一次
func (j *TaskJob) EveryTenMinutes() (cancel CancelFunc) {
	j.interval = time.Minute * 10
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryTwentyMinutes 每20分钟执行一次
func (j *TaskJob) EveryTwentyMinutes() (cancel CancelFunc) {
	j.interval = time.Minute * 20
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryThirtyMinutes 每30分钟执行一次
func (j *TaskJob) EveryThirtyMinutes() (cancel CancelFunc) {
	j.interval = time.Minute * 30
	j.run()

	cancel = j.cancelFunc

	return
}

// Hourly 每1小时执行一次
func (j *TaskJob) Hourly() (cancel CancelFunc) {
	j.interval = time.Hour
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryFiveHours 每5小时执行一次
func (j *TaskJob) EveryFiveHours() (cancel CancelFunc) {
	j.interval = time.Hour * 5
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryTenHours 每10小时执行一次
func (j *TaskJob) EveryTenHours() (cancel CancelFunc) {
	j.interval = time.Hour * 10
	j.run()

	cancel = j.cancelFunc

	return
}

// EveryTwentyHours 每20小时执行一次
func (j *TaskJob) EveryTwentyHours() (cancel CancelFunc) {
	j.interval = time.Hour * 20
	j.run()

	cancel = j.cancelFunc

	return
}

// Daily 每天执行一次
func (j *TaskJob) Daily() (cancel CancelFunc) {
	j.interval = time.Hour * 24
	j.run()

	cancel = j.cancelFunc

	return
}

// Weekly 每周执行一次（每7天）
func (j *TaskJob) Weekly() (cancel CancelFunc) {
	j.interval = time.Hour * 24 * 7
	j.run()

	cancel = j.cancelFunc

	return
}

// Monthly 每月执行一次（每30天）
func (j *TaskJob) Monthly() (cancel CancelFunc) {
	j.interval = time.Hour * 24 * 30
	j.run()

	cancel = j.cancelFunc

	return
}

// Yearly 每年执行一次（每365天）
func (j *TaskJob) Yearly() (cancel CancelFunc) {
	j.interval = time.Hour * 24 * 365
	j.run()

	cancel = j.cancelFunc

	return
}
