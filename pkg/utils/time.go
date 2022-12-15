package utils

import (
	"github.com/keepchen/go-sail/pkg/app/user/config"
	"log"
	"time"
)

type TIM struct {
	loc            *time.Location
	datetimeLayout string
}

func NewTimeWithTimeZone(timeZone ...string) *TIM {
	var tz string
	if len(timeZone) > 0 {
		tz = timeZone[0]
	} else {
		tz = config.GetGlobalConfig().Timezone
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
