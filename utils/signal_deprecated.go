package utils

import (
	"fmt"
	"os/signal"
	"sync"
	"syscall"
)

// ListeningExitSignal 监听系统退出信号
//
// Deprecated: ListeningExitSignal is deprecated,it will be removed in the future.
//
// Please use Signal().ListeningExit() instead.
func ListeningExitSignal(wg *sync.WaitGroup) {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals

	wg.Wait()
	fmt.Printf("Receive signal: %v,program exited.\n", sig)
}
