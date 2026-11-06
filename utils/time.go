package utils

import (
	"log"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
)

type TIM struct {
	loc                       *time.Location
	datetimeTZLayoutWithMilli string
	datetimeTZLayout          string
	datetimeLayout            string
	dateLayout                string
	timeLayout                string
	dirtyTimezone             bool
}

// NewTimeWithTimeZone initializes a time using the specified time zone
//
// Default time zone: Asia/Shanghai
func NewTimeWithTimeZone(timeZone ...string) *TIM {
	var (
		tz            string
		dirtyTimezone bool
	)
	if len(timeZone) > 0 {
		tz = timeZone[0]
	} else {
		tz = constants.DefaultTimeZone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("load location error: %s", err.Error())
		dirtyTimezone = true
	}

	return &TIM{
		loc:                       loc,
		datetimeTZLayoutWithMilli: constants.DateTimeTZLayoutWithMilli,
		datetimeTZLayout:          constants.DatetimeTZLayout,
		datetimeLayout:            constants.DatetimeLayout,
		dateLayout:                constants.DateLayout,
		timeLayout:                constants.TimeLayout,
		dirtyTimezone:             dirtyTimezone,
	}
}

// Now retrieves the current time object (with time zone)
func (t *TIM) Now() time.Time {
	return time.Now().In(t.loc)
}

// DatetimeTZ retrieves the formatted current date and time with time zone
func (t *TIM) DatetimeTZ() string {
	return t.Now().Format(t.datetimeTZLayout)
}

// Datetime retrieves the formatted current date and time
func (t *TIM) Datetime() string {
	return t.Now().Format(t.datetimeLayout)
}

// Date retrieves the formatted current date
func (t *TIM) Date() string {
	return t.Now().Format(t.dateLayout)
}

// Time retrieves the formatted current time
func (t *TIM) Time() string {
	return t.Now().Format(t.timeLayout)
}
