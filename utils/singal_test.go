package utils

import (
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {
	t.Run("Signal", func(t *testing.T) {
		t.Log(Signal())
	})
}

func TestListeningExit(t *testing.T) {
	t.Run("ListeningExit", func(t *testing.T) {
		wg := &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			time.Sleep(time.Second)
			wg.Done()

			//模拟触发
			signals <- syscall.SIGINT
		}()

		Signal().ListeningExit(wg)
	})
}
