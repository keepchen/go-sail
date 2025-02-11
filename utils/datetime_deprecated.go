package utils

import (
	"strings"
	"time"
)

// FormatDate 时间对象转字符串
//
// Deprecated: FormatDate is deprecated,it will be removed in the future.
//
// Please use Datetime().FormatDate() instead.
func FormatDate(date time.Time, dateStyle DateStyle) string {
	layout := string(dateStyle)
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.Replace(layout, "SSS", "000", -1)

	return date.Format(layout)
}

// ParseDate 解析时间
//
// Deprecated: FormatDate is deprecated,it will be removed in the future.
//
// Please use Datetime().ParseDate() instead.
//
// string解析到time对象
func ParseDate(date, layout string, loc *time.Location) (time.Time, error) {
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.Replace(layout, "SSS", "000", -1)

	if loc == nil {
		loc = time.Local
	}

	return time.ParseInLocation(layout, date, loc)
}
