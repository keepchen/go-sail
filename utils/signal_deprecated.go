package utils

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

// ListeningExitSignal 监听系统退出信号
//
// Deprecated: ListeningExitSignal is deprecated,it will be removed in the future.
//
// Please use Signal().ListeningExit() instead.
func ListeningExitSignal(wg *sync.WaitGroup) {
	signals := make(chan os.Signal, 1) // 监听退出
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals

	wg.Wait()
	logger.GetLogger().Sugar().Infof("Receive signal: %v,program exited.", zap.Any("signal", sig))
}
