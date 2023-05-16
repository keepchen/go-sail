package utils

import (
	"log"
	"time"

	"github.com/keepchen/go-sail/pkg/constants"
)

type TIM struct {
	loc            *time.Location
	datetimeLayout string
}

// NewTimeWithTimeZone 根据时区初始化时间
//
// 默认时区: Asia/Shanghai
func NewTimeWithTimeZone(timeZone ...string) *TIM {
	var tz string
	if len(timeZone) > 0 {
		tz = timeZone[0]
	} else {
		tz = constants.DefaultTimeZone
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("load location error: %s", err.Error())
	}

	return &TIM{
		loc:            loc,
		datetimeLayout: "2006-01-02 15:04:05",
	}
}

func (t *TIM) Now() time.Time {
	return time.Now().In(t.loc)
}

func (t *TIM) Datetime() string {
	return t.Now().Format(t.datetimeLayout)
}
