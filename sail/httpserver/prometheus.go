package httpserver

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RunPrometheusServerWhenEnable 启动prometheus指标收集服务
//
// 当配置文件指明启用时才会启动
func RunPrometheusServerWhenEnable(conf config.PrometheusConf) {
	if !conf.Enable {
		return
	}
	var path = conf.AccessPath
	if len(path) == 0 {
		path = "/metrics"
	}

	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())
	srv := &http.Server{
		Addr:    conf.Addr,
		Handler: mux,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Info("Prometheus listen error", zap.String("err", err.Error()))
		}
	}()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.GetLogger().Error("Prometheus Server forced to shutdown", zap.String("err", err.Error()))
		}
	}()
}
