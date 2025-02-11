package utils

import "testing"

func TestRedisLockerImplLockerValue(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(lockerValue())
	}
}
