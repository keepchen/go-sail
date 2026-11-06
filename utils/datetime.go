package utils

import (
	"strings"
	"time"
)

// DateStyle date format
type DateStyle string

const (
	MM_DD                   = "MM-dd"
	YYYYMM                  = "yyyyMM"
	YYYY_MM                 = "yyyy-MM"
	YYYY_MM_DD              = "yyyy-MM-dd"
	YYYYMMDD                = "yyyyMMdd"
	YYYYMMDDHHMMSS          = "yyyyMMddHHmmss"
	YYYYMMDDHHMM            = "yyyyMMddHHmm"
	YYYYMMDDHH              = "yyyyMMddHH"
	YYMMDDHHMM              = "yyMMddHHmm"
	MM_DD_HH_MM             = "MM-dd HH:mm"
	MM_DD_HH_MM_SS          = "MM-dd HH:mm:ss"
	YYYY_MM_DD_HH_MM        = "yyyy-MM-dd HH:mm"
	YYYY_MM_DD_HH_MM_SS     = "yyyy-MM-dd HH:mm:ss"
	YYYY_MM_DD_HH_MM_SS_SSS = "yyyy-MM-dd HH:mm:ss.SSS"

	MM_DD_EN                   = "MM/dd"
	YYYY_MM_EN                 = "yyyy/MM"
	YYYY_MM_DD_EN              = "yyyy/MM/dd"
	MM_DD_HH_MM_EN             = "MM/dd HH:mm"
	MM_DD_HH_MM_SS_EN          = "MM/dd HH:mm:ss"
	YYYY_MM_DD_HH_MM_EN        = "yyyy/MM/dd HH:mm"
	YYYY_MM_DD_HH_MM_SS_EN     = "yyyy/MM/dd HH:mm:ss"
	YYYY_MM_DD_HH_MM_SS_SSS_EN = "yyyy/MM/dd HH:mm:ss.SSS"

	MM_DD_CN               = "MM月dd日"
	YYYY_MM_CN             = "yyyy年MM月"
	YYYY_MM_DD_CN          = "yyyy年MM月dd日"
	MM_DD_HH_MM_CN         = "MM月dd日 HH:mm"
	MM_DD_HH_MM_SS_CN      = "MM月dd日 HH:mm:ss"
	YYYY_MM_DD_HH_MM_CN    = "yyyy年MM月dd日 HH:mm"
	YYYY_MM_DD_HH_MM_SS_CN = "yyyy年MM月dd日 HH:mm:ss"

	HH_MM       = "HH:mm"
	HH_MM_SS    = "HH:mm:ss"
	HH_MM_SS_MS = "HH:mm:ss.SSS"
)

type datetimeImpl struct {
}

type IDatetime interface {
	// FormatDate converts a time object to a string
	FormatDate(date time.Time, dateStyle DateStyle) string
	// ParseDate parses a time value
	//
	// Parse a string into a time object
	ParseDate(date, layout string, loc *time.Location) (time.Time, error)
}

var dti IDatetime = &datetimeImpl{}

// Datetime instantiates date and time utilities
func Datetime() IDatetime {
	return dti
}

// FormatDate converts a time object to a string
func (datetimeImpl) FormatDate(date time.Time, dateStyle DateStyle) string {
	layout := string(dateStyle)
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.ReplaceAll(layout, "SSS", "000")

	return date.Format(layout)
}

// ParseDate parses a time value
//
// Parse a string into a time object
func (datetimeImpl) ParseDate(date, layout string, loc *time.Location) (time.Time, error) {
	layout = strings.Replace(layout, "yyyy", "2006", 1)
	layout = strings.Replace(layout, "yy", "06", 1)
	layout = strings.Replace(layout, "MM", "01", 1)
	layout = strings.Replace(layout, "dd", "02", 1)
	layout = strings.Replace(layout, "HH", "15", 1)
	layout = strings.Replace(layout, "mm", "04", 1)
	layout = strings.Replace(layout, "ss", "05", 1)
	layout = strings.ReplaceAll(layout, "SSS", "000")

	if loc == nil {
		loc = time.Local
	}

	return time.ParseInLocation(layout, date, loc)
}
