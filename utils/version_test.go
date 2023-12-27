package utils

import (
	"runtime"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/constants"
)

func TestPrintVersion(t *testing.T) {
	var fields = VersionInfoFields{
		AppName:   "go-sail",
		Version:   constants.GoSailVersion,
		Branch:    "main",
		Revision:  "cf6e7f1",
		BuildDate: FormatDate(time.Now(), YYYY_MM_DD_HH_MM_SS_EN),
		GoVersion: runtime.Version(),
	}

	PrintVersion(fields)
}
