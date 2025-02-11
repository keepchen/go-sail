package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

type signalImpl struct {
}

type ISignal interface {
	// ListeningExit 监听系统退出信号
	ListeningExit(wg *sync.WaitGroup)
}

var _ ISignal = &signalImpl{}

// Signal 实例化信号工具类
func Signal() ISignal {
	return &signalImpl{}
}

// ListeningExit 监听系统退出信号
func (signalImpl) ListeningExit(wg *sync.WaitGroup) {
	signals := make(chan os.Signal, 1) // 监听退出
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals

	wg.Wait()
	logger.GetLogger().Sugar().Infof("Receive signal: %v,program exited.", zap.Any("signal", sig))
}
