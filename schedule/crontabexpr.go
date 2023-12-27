package schedule

// 一些常用的crontab表达式
const (
	FirstDayOfMonth = "0 0 1 * *" //每月的第一天
	LastDayOfMonth  = "0 0 L * *" //每月的最后一天
	FirstDayOfWeek  = "0 0 * * 1" //每周的第一天（周一）
	LastDayOfWeek   = "0 0 * * 7" //每周的最后一天（周天）
)
