package schedule

// 一些常用的crontab表达式
const (
	EveryMinute                           = "* * * * *"             //每分钟的开始第0秒
	EveryFiveMinute                       = "*/5 * * * *"           //每5分钟的开始第0秒
	EveryTenMinute                        = "*/10 * * * *"          //每10分钟的开始第0秒
	EveryFifteenMinute                    = "*/15 * * * *"          //每15分钟的开始第0秒
	EveryTwentyMinute                     = "*/20 * * * *"          //每20分钟的开始第0秒
	EveryThirtyMinute                     = "*/30 * * * *"          //每30分钟的开始第0秒
	EveryFortyFiveMinute                  = "*/45 * * * *"          //每45分钟的开始第0秒
	FirstDayOfMonth                       = "0 0 1 * *"             //每月的第一天的0点0分
	LastDayOfMonth                        = "0 0 L * *"             //每月的最后一天的0点0分
	FirstDayOfWeek                        = "0 0 * * 1"             //每周的第一天（周一）的0点0分
	LastDayOfWeek                         = "0 0 * * 6"             //每周的最后一天（周天）的0点0分
	DailyAtTenAM                          = "0 10 * * *"            //每天上午10点
	DailyAtTwentyPM                       = "0 20 * * *"            //每天晚上20点
	TenClockAtWeekday                     = "0 10 * * MON-FRI"      //每个工作日（周一~周五）的上午10点0分
	TenClockAtWeekend                     = "0 10 * * SAT,SUN"      //每个周末（周六和周日）的上午10点0分
	HourlyBetween9And17ClockAtWeekday     = "0 9-17 * * MON-FRI"    //每个工作日（周一~周五）的上午9点0分到下午5点0分每小时一次
	HalfHourlyBetween9And17ClockAtWeekday = "*/30 9-17 * * MON-FRI" //每个工作日（周一~周五）的上午9点0分到下午5点0分每半时一次
)
