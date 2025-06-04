package schedule

import "testing"

func TestPrintExpr(t *testing.T) {
	t.Run("PrintExpr", func(t *testing.T) {
		t.Log(EveryMinute)
		t.Log(EveryFiveMinute)
		t.Log(EveryTenMinute)
		t.Log(EveryFifteenMinute)
		t.Log(EveryTwentyMinute)
		t.Log(EveryThirtyMinute)
		t.Log(EveryFortyFiveMinute)
		t.Log(FirstDayOfMonth)
		t.Log(LastDayOfMonth)
		t.Log(FirstDayOfWeek)
		t.Log(LastDayOfWeek)
		t.Log(DailyAtTenAM)
		t.Log(DailyAtTwentyPM)
		t.Log(TenClockAtWeekday)
		t.Log(TenClockAtWeekend)
		t.Log(HourlyBetween9And17ClockAtWeekday)
		t.Log(HalfHourlyBetween9And17ClockAtWeekday)
	})
}
