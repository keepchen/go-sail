package utils

import (
	"log"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
)

type TIM struct {
	loc            *time.Location
	datetimeLayout string
	dateLayout     string
	timeLayout     string
	dirtyTimezone  bool
}

// NewTimeWithTimeZone 根据时区初始化时间
//
// 默认时区: Asia/Shanghai
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
		loc:            loc,
		datetimeLayout: "2006-01-02 15:04:05",
		dateLayout:     "2006-01-02",
		timeLayout:     "15:04:05",
		dirtyTimezone:  dirtyTimezone,
	}
}

func (t *TIM) Now() time.Time {
	return time.Now().In(t.loc)
}

func (t *TIM) Datetime() string {
	return t.Now().Format(t.datetimeLayout)
}

func (t *TIM) Date() string {
	return t.Now().Format(t.dateLayout)
}

func (t *TIM) Time() string {
	return t.Now().Format(t.timeLayout)
}
