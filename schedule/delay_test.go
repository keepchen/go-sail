package schedule

import (
	"fmt"
	"testing"
	"time"
)

func TestRunAfter(t *testing.T) {
	t.Run("RunAfter", func(t *testing.T) {
		NewJob("RunAfter", func() {
			fmt.Println("RunAfter...")
		}).RunAfter(time.Second)

		cancel := NewJob("RunAfter2", func() {
			fmt.Println("RunAfter2...")
		}).RunAfter(time.Minute)

		cancel()
	})
}
