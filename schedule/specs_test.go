package schedule

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
)

func TestCronJobRegistration(t *testing.T) {
	tests := []struct {
		name     string
		register func(Scheduler) CancelFunc
		want     string
	}{
		{"RunAtFiveFields", func(j Scheduler) CancelFunc { return j.RunAt("* * * * *") }, "2026-11-02T09:01:00Z"},
		{"RunAtSixFields", func(j Scheduler) CancelFunc { return j.RunAt(EveryMinute) }, "2026-11-02T09:01:00Z"},
		{"RunAtSeconds", func(j Scheduler) CancelFunc { return j.RunAt("*/5 * * * * *") }, "2026-11-02T09:00:05Z"},
		{"RunAtMonthEnd", func(j Scheduler) CancelFunc { return j.RunAt(LastDayOfMonth) }, "2026-11-30T00:00:00Z"},
		{"TenClockAtWeekday", func(j Scheduler) CancelFunc { return j.TenClockAtWeekday() }, "2026-11-02T10:00:00Z"},
		{"TenClockAtWeekend", func(j Scheduler) CancelFunc { return j.TenClockAtWeekend() }, "2026-11-07T10:00:00Z"},
		{"FirstDayOfMonthly", func(j Scheduler) CancelFunc { return j.FirstDayOfMonthly() }, "2026-12-01T00:00:00Z"},
		{"LastDayOfMonthly", func(j Scheduler) CancelFunc { return j.LastDayOfMonthly() }, "2026-11-30T00:00:00Z"},
		{"FirstDayOfWeek", func(j Scheduler) CancelFunc { return j.FirstDayOfWeek() }, "2026-11-09T00:00:00Z"},
		{"LastDayOfWeek", func(j Scheduler) CancelFunc { return j.LastDayOfWeek() }, "2026-11-08T00:00:00Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := make(map[cron.EntryID]bool)
			if cronJob != nil {
				for _, entry := range cronJob.Entries() {
					before[entry.ID] = true
				}
			}
			var called atomic.Bool
			job := NewJob(t.Name(), func() { called.Store(true) })
			cancel := tt.register(job)
			t.Cleanup(func() {
				cancel()
				// Cancellation is asynchronous. Wait for it so repeated test runs do not
				// leave duplicate names or pending entries in the shared scheduler.
				deadline := time.Now().Add(time.Second)
				for time.Now().Before(deadline) {
					taskSchedules.mux.RLock()
					_, exists := taskSchedules.pool[generateJobNameKey(t.Name())]
					taskSchedules.mux.RUnlock()
					if !exists {
						return
					}
					time.Sleep(time.Millisecond)
				}
				t.Error("job was not removed after cancellation")
			})
			var added []cron.Entry
			for _, entry := range cronJob.Entries() {
				if !before[entry.ID] {
					added = append(added, entry)
				}
			}
			if len(added) != 1 {
				t.Fatalf("registered %d entries, want 1", len(added))
			}
			entry := added[0]
			if entry.Next.IsZero() {
				t.Fatal("registered job has no next execution")
			}
			got := entry.Schedule.Next(parseTestTime(t, "2026-11-02T09:00:01Z"))
			if want := parseTestTime(t, tt.want); !got.Equal(want) {
				t.Fatalf("next = %s, want %s", got, want)
			}
			entry.Job.Run()
			if !called.Load() {
				t.Error("registered callback did not execute the task")
			}
		})
	}
}
