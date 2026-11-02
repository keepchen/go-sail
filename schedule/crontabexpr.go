package schedule

// 一些常用的crontab表达式
const (
	EveryMinute                           = "0 * * * * *"             //每分钟的开始第0秒
	EveryFiveMinute                       = "0 */5 * * * *"           //每5分钟的开始第0秒
	EveryTenMinute                        = "0 */10 * * * *"          //每10分钟的开始第0秒
	EveryFifteenMinute                    = "0 */15 * * * *"          //每15分钟的开始第0秒
	EveryTwentyMinute                     = "0 */20 * * * *"          //每20分钟的开始第0秒
	EveryThirtyMinute                     = "0 */30 * * * *"          //每30分钟的开始第0秒
	EveryFortyFiveMinute                  = "0 */45 * * * *"          //每小时的第0和第45分钟的第0秒
	FirstDayOfMonth                       = "0 0 0 1 * *"             //每月的第一天的0点0分
	LastDayOfMonth                        = "0 0 0 L * *"             //每月的最后一天的0点0分（由go-sail解析器支持）
	FirstDayOfWeek                        = "0 0 0 * * 1"             //每周的第一天（周一）的0点0分
	LastDayOfWeek                         = "0 0 0 * * SUN"           //每周的最后一天（周天）的0点0分
	DailyAtTenAM                          = "0 0 10 * * *"            //每天上午10点
	DailyAtTwentyPM                       = "0 0 20 * * *"            //每天晚上20点
	TenClockAtWeekday                     = "0 0 10 * * MON-FRI"      //每个工作日（周一~周五）的上午10点0分
	TenClockAtWeekend                     = "0 0 10 * * SAT,SUN"      //每个周末（周六和周日）的上午10点0分
	HourlyBetween9And17ClockAtWeekday     = "0 0 9-17 * * MON-FRI"    //每个工作日（周一~周五）的上午9点0分到下午5点0分每小时一次
	HalfHourlyBetween9And17ClockAtWeekday = "0 */30 9-17 * * MON-FRI" //每个工作日（周一~周五）的上午9点0分到下午5点30分每半小时一次
)
