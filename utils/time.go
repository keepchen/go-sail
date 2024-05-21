package utils

import (
	"log"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
)

type TIM struct {
	loc              *time.Location
	datetimeTZLayout string
	datetimeLayout   string
	dateLayout       string
	timeLayout       string
	dirtyTimezone    bool
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
		loc:              loc,
		datetimeTZLayout: constants.DatetimeTZLayout,
		datetimeLayout:   constants.DatetimeLayout,
		dateLayout:       constants.DateLayout,
		timeLayout:       constants.TimeLayout,
		dirtyTimezone:    dirtyTimezone,
	}
}

// Now 获取当前时间对象（带时区）
func (t *TIM) Now() time.Time {
	return time.Now().In(t.loc)
}

// DatetimeTZ 获取格式化后的当前日期时间-tz
func (t *TIM) DatetimeTZ() string {
	return t.Now().Format(t.datetimeTZLayout)
}

// Datetime 获取格式化后的当前日期时间
func (t *TIM) Datetime() string {
	return t.Now().Format(t.datetimeLayout)
}

// Date 获取格式化后的当前日期
func (t *TIM) Date() string {
	return t.Now().Format(t.dateLayout)
}

// Time 获取格式化后的当前时间
func (t *TIM) Time() string {
	return t.Now().Format(t.timeLayout)
}
