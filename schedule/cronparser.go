package schedule

import (
	"time"

	"github.com/robfig/cron/v3"
)

// scheduleParser accepts optional seconds and adds support for our month-end preset.
// Other Quartz extensions are not supported.
type scheduleParser struct{}

func (scheduleParser) Parse(spec string) (cron.Schedule, error) {
	if spec == LastDayOfMonth {
		return monthEndSchedule{}, nil
	}
	return cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom |
			cron.Month | cron.Dow | cron.Descriptor,
	).Parse(spec)
}

type monthEndSchedule struct{}

func (monthEndSchedule) Next(after time.Time) time.Time {
	// Day zero of the next month is the last day of the current month.
	next := time.Date(after.Year(), after.Month()+1, 0, 0, 0, 0, 0, after.Location())
	if !next.After(after) {
		next = time.Date(after.Year(), after.Month()+2, 0, 0, 0, 0, 0, after.Location())
	}
	return next
}
