package orm

import "testing"

func TestNowTimeFunc(t *testing.T) {
	t.Run("NowTimeFunc", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Log(nowTimeFunc())
		}
	})
}

func TestSetHookTimeFunc(t *testing.T) {
	t.Run("SetHookTimeFunc", func(t *testing.T) {
		SetHookTimeFunc(nowTimeFunc)
	})
}
