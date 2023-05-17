package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/pkg/constants"
)

var timezones = []string{
	constants.DefaultTimeZone,
	constants.TimeZoneUTCSub11,
	constants.TimeZoneUTCSub10,
	constants.TimeZoneUTCSub9,
	constants.TimeZoneUTCSub8,
	constants.TimeZoneUTCSub7,
	constants.TimeZoneUTCSub6,
	constants.TimeZoneUTCSub5,
	constants.TimeZoneUTCSub4,
	constants.TimeZoneUTCSub3,
	constants.TimeZoneUTCSub2,
	constants.TimeZoneUTCSub1,
	constants.TimeZoneUTC0,
	constants.TimeZoneUTCPlus1,
	constants.TimeZoneUTCPlus2,
	constants.TimeZoneUTCPlus3,
	constants.TimeZoneUTCPlus4,
	constants.TimeZoneUTCPlus5,
	constants.TimeZoneUTCPlus6,
	constants.TimeZoneUTCPlus7,
	constants.TimeZoneUTCPlus8,
	constants.TimeZoneUTCPlus9,
	constants.TimeZoneUTCPlus10,
	constants.TimeZoneUTCPlus11,
	constants.TimeZoneUTCPlus12,
	"Earth/Unknown",
}

func TestNewTimeWithTimeZone(t *testing.T) {
	for index, timezone := range timezones {
		tim := NewTimeWithTimeZone(timezone)
		if index == len(timezones)-1 {
			assert.Equal(t, true, tim.dirtyTimezone)
		} else {
			assert.Equal(t, false, tim.dirtyTimezone)
		}
	}
}
