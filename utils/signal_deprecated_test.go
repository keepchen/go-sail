package utils

import (
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestListeningExitSignal(t *testing.T) {
	t.Run("ListeningExitSignal", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			time.Sleep(time.Second)
			wg.Done()

			//模拟触发
			signals <- syscall.SIGINT
		}()

		ListeningExitSignal(wg)
	})
}
