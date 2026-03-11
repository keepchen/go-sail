package constants

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrintTimezone(t *testing.T) {
	zones := []string{
		DefaultTimeZone,
		TimeZoneUTCSub11,
		TimeZoneUTCSub10,
		TimeZoneUTCSub9,
		TimeZoneUTCSub8,
		TimeZoneUTCSub7,
		TimeZoneUTCSub6,
		TimeZoneUTCSub5,
		TimeZoneUTCSub4,
		TimeZoneUTCSub3,
		TimeZoneUTCSub2,
		TimeZoneUTCSub1,
		TimeZoneUTC0,
		TimeZoneUTCPlus1,
		TimeZoneUTCPlus2,
		TimeZoneUTCPlus3,
		TimeZoneUTCPlus4,
		TimeZoneUTCPlus5,
		TimeZoneUTCPlus6,
		TimeZoneUTCPlus7,
		TimeZoneUTCPlus8,
		TimeZoneUTCPlus9,
		TimeZoneUTCPlus10,
		TimeZoneUTCPlus11,
		TimeZoneUTCPlus12,
	}
	for _, zone := range zones {
		loc, err := time.LoadLocation(zone)
		assert.NoError(t, err)
		assert.NotNil(t, loc)
	}
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
