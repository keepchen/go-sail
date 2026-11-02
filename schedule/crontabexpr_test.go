package schedule

import (
	"testing"
	"time"
)

func parseTestTime(t *testing.T, value string) time.Time {
	t.Helper()
	result, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestCrontabExpressions(t *testing.T) {
	// Monday, just after the hour: verifies that seconds are zero and the next
	// occurrence is strictly later than the reference time.
	after := parseTestTime(t, "2026-11-02T09:00:01Z")
	tests := []struct {
		name, expr string
		want       []string
	}{
		{"EveryMinute", EveryMinute, []string{"2026-11-02T09:01:00Z", "2026-11-02T09:02:00Z", "2026-11-02T09:03:00Z"}},
		{"EveryFiveMinute", EveryFiveMinute, []string{"2026-11-02T09:05:00Z", "2026-11-02T09:10:00Z", "2026-11-02T09:15:00Z"}},
		{"EveryTenMinute", EveryTenMinute, []string{"2026-11-02T09:10:00Z", "2026-11-02T09:20:00Z", "2026-11-02T09:30:00Z"}},
		{"EveryFifteenMinute", EveryFifteenMinute, []string{"2026-11-02T09:15:00Z", "2026-11-02T09:30:00Z", "2026-11-02T09:45:00Z"}},
		{"EveryTwentyMinute", EveryTwentyMinute, []string{"2026-11-02T09:20:00Z", "2026-11-02T09:40:00Z", "2026-11-02T10:00:00Z"}},
		{"EveryThirtyMinute", EveryThirtyMinute, []string{"2026-11-02T09:30:00Z", "2026-11-02T10:00:00Z", "2026-11-02T10:30:00Z"}},
		{"EveryFortyFiveMinute", EveryFortyFiveMinute, []string{"2026-11-02T09:45:00Z", "2026-11-02T10:00:00Z", "2026-11-02T10:45:00Z"}},
		{"FirstDayOfMonth", FirstDayOfMonth, []string{"2026-12-01T00:00:00Z", "2027-01-01T00:00:00Z", "2027-02-01T00:00:00Z"}},
		{"LastDayOfMonth", LastDayOfMonth, []string{"2026-11-30T00:00:00Z", "2026-12-31T00:00:00Z", "2027-01-31T00:00:00Z"}},
		{"FirstDayOfWeek", FirstDayOfWeek, []string{"2026-11-09T00:00:00Z", "2026-11-16T00:00:00Z", "2026-11-23T00:00:00Z"}},
		{"LastDayOfWeek", LastDayOfWeek, []string{"2026-11-08T00:00:00Z", "2026-11-15T00:00:00Z", "2026-11-22T00:00:00Z"}},
		{"DailyAtTenAM", DailyAtTenAM, []string{"2026-11-02T10:00:00Z", "2026-11-03T10:00:00Z", "2026-11-04T10:00:00Z"}},
		{"DailyAtTwentyPM", DailyAtTwentyPM, []string{"2026-11-02T20:00:00Z", "2026-11-03T20:00:00Z", "2026-11-04T20:00:00Z"}},
		{"TenClockAtWeekday", TenClockAtWeekday, []string{"2026-11-02T10:00:00Z", "2026-11-03T10:00:00Z", "2026-11-04T10:00:00Z"}},
		{"TenClockAtWeekend", TenClockAtWeekend, []string{"2026-11-07T10:00:00Z", "2026-11-08T10:00:00Z", "2026-11-14T10:00:00Z"}},
		{"HourlyBetween9And17ClockAtWeekday", HourlyBetween9And17ClockAtWeekday, []string{"2026-11-02T10:00:00Z", "2026-11-02T11:00:00Z", "2026-11-02T12:00:00Z"}},
		{"HalfHourlyBetween9And17ClockAtWeekday", HalfHourlyBetween9And17ClockAtWeekday, []string{"2026-11-02T09:30:00Z", "2026-11-02T10:00:00Z", "2026-11-02T10:30:00Z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := (scheduleParser{}).Parse(tt.expr)
			if err != nil {
				t.Fatalf("parse %q: %v", tt.expr, err)
			}
			next := after
			for _, expected := range tt.want {
				next = schedule.Next(next)
				t.Logf("%s (%s) -> %s", tt.name, tt.expr, next.Format(time.RFC3339))
				if want := parseTestTime(t, expected); !next.Equal(want) {
					t.Fatalf("next = %s, want %s", next, want)
				}
			}
		})
	}
}

func TestScheduleParserCompatibility(t *testing.T) {
	for _, tt := range []struct{ expr, want string }{
		{"* * * * *", "2026-11-02T09:01:00Z"},
		{"0 * * * * *", "2026-11-02T09:01:00Z"},
		{"*/5 * * * * *", "2026-11-02T09:00:05Z"},
		{"@daily", "2026-11-03T00:00:00Z"},
		{"@every 5s", "2026-11-02T09:00:06Z"},
		{"CRON_TZ=Asia/Shanghai 0 0 10 * * *", "2026-11-03T02:00:00Z"},
	} {
		t.Run(tt.expr, func(t *testing.T) {
			schedule, err := (scheduleParser{}).Parse(tt.expr)
			if err != nil {
				t.Fatal(err)
			}
			got := schedule.Next(parseTestTime(t, "2026-11-02T09:00:01Z"))
			if want := parseTestTime(t, tt.want); !got.Equal(want) {
				t.Fatalf("next = %s, want %s", got, want)
			}
		})
	}
	for _, expr := range []string{"", "* * * *", "* * * * * * *", "60 * * * * *", "0 0 0 L 2 *"} {
		if _, err := (scheduleParser{}).Parse(expr); err == nil {
			t.Errorf("expected parse error for %q", expr)
		}
	}
}

func TestMonthEndSchedule(t *testing.T) {
	for _, tt := range []struct{ name, after, want string }{
		{"common February", "2027-02-01T00:00:00Z", "2027-02-28T00:00:00Z"},
		{"leap February", "2028-02-01T00:00:00Z", "2028-02-29T00:00:00Z"},
		{"century common year", "2100-02-01T00:00:00Z", "2100-02-28T00:00:00Z"},
		{"century leap year", "2000-02-01T00:00:00Z", "2000-02-29T00:00:00Z"},
		{"30 days", "2026-04-01T00:00:00Z", "2026-04-30T00:00:00Z"},
		{"31 days", "2026-05-01T00:00:00Z", "2026-05-31T00:00:00Z"},
		{"before midnight", "2026-04-29T23:59:59Z", "2026-04-30T00:00:00Z"},
		{"at midnight", "2026-04-30T00:00:00Z", "2026-05-31T00:00:00Z"},
		{"after midnight", "2026-04-30T00:00:01Z", "2026-05-31T00:00:00Z"},
		{"year rollover", "2026-12-31T00:00:00Z", "2027-01-31T00:00:00Z"},
		{"local midnight", "2028-02-28T23:59:59+08:00", "2028-02-29T00:00:00+08:00"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			after := parseTestTime(t, tt.after)
			got := (monthEndSchedule{}).Next(after)
			want := parseTestTime(t, tt.want)
			t.Logf("%s -> %s", after.Format(time.RFC3339), got.Format(time.RFC3339))
			if !got.Equal(want) || got.Location() != after.Location() {
				t.Fatalf("next = %s, want %s in original location", got, want)
			}
		})
	}
}

func TestWeekdayHourBoundaries(t *testing.T) {
	for _, tt := range []struct{ name, expr, after, want string }{
		{"weekday skips weekend", TenClockAtWeekday, "2026-11-06T10:00:00Z", "2026-11-09T10:00:00Z"},
		{"hourly end of workday", HourlyBetween9And17ClockAtWeekday, "2026-11-06T17:00:00Z", "2026-11-09T09:00:00Z"},
		{"half hourly includes 1730", HalfHourlyBetween9And17ClockAtWeekday, "2026-11-06T17:00:00Z", "2026-11-06T17:30:00Z"},
		{"half hourly skips weekend", HalfHourlyBetween9And17ClockAtWeekday, "2026-11-06T17:30:00Z", "2026-11-09T09:00:00Z"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			schedule, err := (scheduleParser{}).Parse(tt.expr)
			if err != nil {
				t.Fatal(err)
			}
			got := schedule.Next(parseTestTime(t, tt.after))
			if want := parseTestTime(t, tt.want); !got.Equal(want) {
				t.Fatalf("next = %s, want %s", got, want)
			}
		})
	}
}
