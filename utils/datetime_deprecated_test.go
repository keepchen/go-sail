package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeParse(t *testing.T) {
	datetime, err := ParseDate("2023-12-27", "yyyy-MM-dd", nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, "2023-12-27", datetime.Format("2006-01-02"))

	loc, locErr := time.LoadLocation("Asia/Shanghai")
	assert.Equal(t, nil, locErr)
	datetime, err = ParseDate("2023-12-27 12:12:12", "yyyy-MM-dd HH:mm:ss", loc)
	assert.Equal(t, nil, err)
	assert.Equal(t, "2023-12-27", datetime.Format("2006-01-02"))
}

func TestFormatDate(t *testing.T) {
	var styles = []DateStyle{
		MM_DD,
		YYYYMM,
		YYYY_MM,
		YYYY_MM_DD,
		YYYYMMDD,
		YYYYMMDDHHMMSS,
		YYYYMMDDHHMM,
		YYYYMMDDHH,
		YYMMDDHHMM,
		MM_DD_HH_MM,
		MM_DD_HH_MM_SS,
		YYYY_MM_DD_HH_MM,
		YYYY_MM_DD_HH_MM_SS,
		YYYY_MM_DD_HH_MM_SS_SSS,
		MM_DD_EN,
		YYYY_MM_EN,
		YYYY_MM_DD_EN,
		MM_DD_HH_MM_EN,
		MM_DD_HH_MM_SS_EN,
		YYYY_MM_DD_HH_MM_EN,
		YYYY_MM_DD_HH_MM_SS_EN,
		YYYY_MM_DD_HH_MM_SS_SSS_EN,
		MM_DD_CN,
		YYYY_MM_CN,
		YYYY_MM_DD_CN,
		MM_DD_HH_MM_CN,
		MM_DD_HH_MM_SS_CN,
		YYYY_MM_DD_HH_MM_CN,
		YYYY_MM_DD_HH_MM_SS_CN,
		HH_MM,
		HH_MM_SS,
		HH_MM_SS_MS,
	}

	now := NewTimeWithTimeZone("Asia/Shanghai").Now()
	for _, style := range styles {
		t.Log(FormatDate(now, style))
	}
}
