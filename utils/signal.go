package utils

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type signalImpl struct {
}

type ISignal interface {
	// ListeningExit 监听系统退出信号
	ListeningExit(wg *sync.WaitGroup)
}

var si ISignal = &signalImpl{}

// Signal 实例化信号工具类
func Signal() ISignal {
	return si
}

var signals = make(chan os.Signal, 1)

// ListeningExit 监听系统退出信号
func (signalImpl) ListeningExit(wg *sync.WaitGroup) {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals

	wg.Wait()
	fmt.Printf("Receive signal: %v,program exited.\n", sig)
}
