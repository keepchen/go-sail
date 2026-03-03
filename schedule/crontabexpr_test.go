package schedule

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintExpr(t *testing.T) {
	t.Run("PrintExpr", func(t *testing.T) {
		assert.NotEmpty(t, EveryMinute)
		assert.NotEmpty(t, EveryFiveMinute)
		assert.NotEmpty(t, EveryTenMinute)
		assert.NotEmpty(t, EveryFifteenMinute)
		assert.NotEmpty(t, EveryTwentyMinute)
		assert.NotEmpty(t, EveryThirtyMinute)
		assert.NotEmpty(t, EveryFortyFiveMinute)
		assert.NotEmpty(t, FirstDayOfMonth)
		assert.NotEmpty(t, LastDayOfMonth)
		assert.NotEmpty(t, FirstDayOfWeek)
		assert.NotEmpty(t, LastDayOfWeek)
		assert.NotEmpty(t, DailyAtTenAM)
		assert.NotEmpty(t, DailyAtTwentyPM)
		assert.NotEmpty(t, TenClockAtWeekday)
		assert.NotEmpty(t, TenClockAtWeekend)
		assert.NotEmpty(t, HourlyBetween9And17ClockAtWeekday)
		assert.NotEmpty(t, HalfHourlyBetween9And17ClockAtWeekday)
	})
}
