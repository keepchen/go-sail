package schedule

import (
	"fmt"
	"testing"
)

func TestRunAt(t *testing.T) {
	t.Run("RunAt", func(t *testing.T) {
		cancel := NewJob("RunAt", func() {
			fmt.Println("RunAt...")
		}).RunAt(EveryMinute)

		cancel()
	})
}

func TestTenClockAtWeekday(t *testing.T) {
	t.Run("TenClockAtWeekday", func(t *testing.T) {
		cancel := NewJob("TenClockAtWeekday", func() {
			fmt.Println("TenClockAtWeekday...")
		}).TenClockAtWeekday()

		cancel()
	})
}

func TestTenClockAtWeekend(t *testing.T) {
	t.Run("TenClockAtWeekend", func(t *testing.T) {
		cancel := NewJob("TenClockAtWeekend", func() {
			fmt.Println("TenClockAtWeekend...")
		}).TenClockAtWeekend()

		cancel()
	})
}

func TestFirstDayOfMonthly(t *testing.T) {
	t.Run("FirstDayOfMonthly", func(t *testing.T) {
		cancel := NewJob("FirstDayOfMonthly", func() {
			fmt.Println("FirstDayOfMonthly...")
		}).FirstDayOfMonthly()

		cancel()
	})
}

func TestLastDayOfMonthly(t *testing.T) {
	t.Run("LastDayOfMonthly", func(t *testing.T) {
		cancel := NewJob("LastDayOfMonthly", func() {
			fmt.Println("LastDayOfMonthly...")
		}).LastDayOfMonthly()

		cancel()
	})
}

func TestFirstDayOfWeek(t *testing.T) {
	t.Run("FirstDayOfWeek", func(t *testing.T) {
		cancel := NewJob("FirstDayOfWeek", func() {
			fmt.Println("FirstDayOfWeek...")
		}).FirstDayOfWeek()

		cancel()
	})
}

func TestLastDayOfWeek(t *testing.T) {
	t.Run("LastDayOfWeek", func(t *testing.T) {
		cancel := NewJob("LastDayOfWeek", func() {
			fmt.Println("LastDayOfWeek...")
		}).LastDayOfWeek()

		cancel()
	})
}
