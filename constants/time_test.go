package constants

import (
	"testing"
	"time"
)

func TestPrintTimezone(t *testing.T) {
	t.Log(DefaultTimeZone)
	t.Log(TimeZoneUTCSub11)
	t.Log(TimeZoneUTCSub10)
	t.Log(TimeZoneUTCSub9)
	t.Log(TimeZoneUTCSub8)
	t.Log(TimeZoneUTCSub7)
	t.Log(TimeZoneUTCSub6)
	t.Log(TimeZoneUTCSub5)
	t.Log(TimeZoneUTCSub4)
	t.Log(TimeZoneUTCSub3)
	t.Log(TimeZoneUTCSub2)
	t.Log(TimeZoneUTCSub1)
	t.Log(TimeZoneUTC0)
	t.Log(TimeZoneUTCPlus1)
	t.Log(TimeZoneUTCPlus2)
	t.Log(TimeZoneUTCPlus3)
	t.Log(TimeZoneUTCPlus4)
	t.Log(TimeZoneUTCPlus5)
	t.Log(TimeZoneUTCPlus6)
	t.Log(TimeZoneUTCPlus7)
	t.Log(TimeZoneUTCPlus8)
	t.Log(TimeZoneUTCPlus9)
	t.Log(TimeZoneUTCPlus10)
	t.Log(TimeZoneUTCPlus11)
	t.Log(TimeZoneUTCPlus12)
}

func TestPrintLayout(t *testing.T) {
	t.Run("printLayout", func(t *testing.T) {
		t.Log(DateLayout)
		t.Log(TimeLayout)
		t.Log(DatetimeLayout)
		t.Log(DatetimeTZLayout)
		t.Log(DateTimeTZLayoutWithMilli)
	})
	t.Run("formatByLayout", func(t *testing.T) {
		now := time.Now()
		t.Log(now.Format(DateLayout))
		t.Log(now.Format(TimeLayout))
		t.Log(now.Format(DatetimeLayout))
		t.Log(now.Format(DatetimeTZLayout))
		t.Log(DateTimeTZLayoutWithMilli)
	})
}
