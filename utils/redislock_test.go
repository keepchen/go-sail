package utils

import "testing"

func TestLockerValue(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(lockerValue())
	}
}
