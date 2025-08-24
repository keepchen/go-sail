package sail

import (
	"fmt"
	"testing"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

func TestSetScheduleRedisClient(t *testing.T) {
	t.Run("SetScheduleRedisClient", func(t *testing.T) {
		SetScheduleRedisClient(redis.GetInstance())
		SetScheduleRedisClient(redis.GetInstance())
		SetScheduleRedisClient(redis.GetClusterInstance())
		SetScheduleRedisClient(redis.GetClusterInstance())
	})
}

func TestSchedule(t *testing.T) {
	t.Run("Schedule", func(t *testing.T) {
		t.Log(Schedule("test-sail-schedule", func() {
			fmt.Println("test-sail-schedule")
		}))
	})
}
