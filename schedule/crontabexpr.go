package schedule

// Some common crontab expressions
const (
	EveryMinute                           = "0 * * * * *"             //At second 0 of every minute
	EveryFiveMinute                       = "0 */5 * * * *"           //At second 0 of every 5th minute
	EveryTenMinute                        = "0 */10 * * * *"          //At second 0 of every 10th minute
	EveryFifteenMinute                    = "0 */15 * * * *"          //At second 0 of every 15th minute
	EveryTwentyMinute                     = "0 */20 * * * *"          //At second 0 of every 20th minute
	EveryThirtyMinute                     = "0 */30 * * * *"          //At second 0 of every 30th minute
	EveryFortyFiveMinute                  = "0 */45 * * * *"          //At second 0 of minutes 0 and 45 of every hour
	FirstDayOfMonth                       = "0 0 0 1 * *"             //At 00:00 on the first day of every month
	LastDayOfMonth                        = "0 0 0 L * *"             //At 00:00 on the last day of every month (supported by the go-sail parser)
	FirstDayOfWeek                        = "0 0 0 * * 1"             //At 00:00 on the first day of every week (Monday)
	LastDayOfWeek                         = "0 0 0 * * SUN"           //At 00:00 on the last day of every week (Sunday)
	DailyAtTenAM                          = "0 0 10 * * *"            //At 10:00 every day
	DailyAtTwentyPM                       = "0 0 20 * * *"            //At 20:00 every day
	TenClockAtWeekday                     = "0 0 10 * * MON-FRI"      //At 10:00 every weekday (Monday-Friday)
	TenClockAtWeekend                     = "0 0 10 * * SAT,SUN"      //At 10:00 every weekend day (Saturday and Sunday)
	HourlyBetween9And17ClockAtWeekday     = "0 0 9-17 * * MON-FRI"    //Hourly from 09:00 to 17:00 every weekday (Monday-Friday)
	HalfHourlyBetween9And17ClockAtWeekday = "0 */30 9-17 * * MON-FRI" //Every half hour from 09:00 to 17:30 every weekday (Monday-Friday)
)
