package schedule

// TenClockAtWeekday 每个工作日（周一~周五）上午10点
func (j *TaskJob) TenClockAtWeekday() (cancel CancelFunc) {
	return j.RunAt(TenClockAtWeekday)
}

// TenClockAtWeekend 每个周末（周六和周日）上午10点
func (j *TaskJob) TenClockAtWeekend() (cancel CancelFunc) {
	return j.RunAt(TenClockAtWeekend)
}

// FirstDayOfMonthly 每月1号
func (j *TaskJob) FirstDayOfMonthly() (cancel CancelFunc) {
	return j.RunAt(FirstDayOfMonth)
}

// LastDayOfMonthly 每月最后一天
func (j *TaskJob) LastDayOfMonthly() (cancel CancelFunc) {
	return j.RunAt(LastDayOfMonth)
}