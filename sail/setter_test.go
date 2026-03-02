package sail

import "testing"

func TestSetRedisClient(t *testing.T) {
	t.Run("SetRedisClient", func(t *testing.T) {
		t.Log(SetRedisClient(nil))
	})
}

func TestSetRedisClientForSchedule(t *testing.T) {
	t.Run("SetRedisClient-ForSchedule", func(t *testing.T) {
		s := SetRedisClient(nil)
		t.Log(s)
		s.ForSchedule()
	})
}

func TestSetRedisClientForRedisLocker(t *testing.T) {
	t.Run("SetRedisClient-ForRedisLocker", func(t *testing.T) {
		s := SetRedisClient(nil)
		t.Log(s)
		s.ForRedisLocker()
	})
}
