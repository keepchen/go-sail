package sail

import (
	"fmt"
	"testing"
)

func TestSchedule(t *testing.T) {
	t.Run("Schedule", func(t *testing.T) {
		t.Log(Schedule("test-sail-schedule", func() {
			fmt.Println("test-sail-schedule")
		}))
	})
}
