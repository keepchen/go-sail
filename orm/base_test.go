package orm

import "testing"

func TestNowTimeFunc(t *testing.T) {
	t.Run("TestNowTimeFunc", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Log(nowTimeFunc())
		}
	})
}
