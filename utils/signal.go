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
	// ListeningExit listens for system exit signals
	ListeningExit(wg *sync.WaitGroup)
}

var si ISignal = &signalImpl{}

// Signal instantiates signal utilities
func Signal() ISignal {
	return si
}

var signals = make(chan os.Signal, 1)

// ListeningExit listens for system exit signals
func (signalImpl) ListeningExit(wg *sync.WaitGroup) {
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	sig := <-signals

	wg.Wait()
	fmt.Printf("Receive signal: %v,program exited.\n", sig)
}
