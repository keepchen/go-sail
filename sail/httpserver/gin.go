package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/keepchen/go-sail/v3/http/middleware"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"
)

var (
	sigChan = make(chan os.Signal, 1)
	errChan = make(chan error)
)

// NotifyExit 通知退出
func NotifyExit(sig os.Signal) {
	sigChan <- sig
	errChan <- fmt.Errorf("%v", <-sigChan)
}

// InitGinEngine 初始化gin引擎
//
// # Note:
//
// 该方法会默认使用 middleware.LogTrace 中间件
func InitGinEngine(conf config.HttpServerConf) *gin.Engine {
	var r *gin.Engine
	if conf.Debug {
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		r = gin.New()
	}

	if len(conf.TrustedProxies) > 0 {
		_ = r.SetTrustedProxies(conf.TrustedProxies)
	}

	r.Use(middleware.LogTrace())
	if conf.Prometheus.Enable {
		r.Use(middleware.PrometheusExporter())
		if !conf.Prometheus.DisableSystemSample {
			middleware.SetDiskPath(conf.Prometheus.DiskPath)
			middleware.SetSampleInterval(conf.Prometheus.SampleInterval)
			go middleware.SystemMetricsSample()
		}
	}

	return r
}

// RunHttpServer 启动http服务
func RunHttpServer(conf config.HttpServerConf, ginEngine *gin.Engine, apiOption *api.Option, wg *sync.WaitGroup) {
	defer wg.Done()

	//监听退出信号
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%v", <-sigChan)
		cancel()
	}()

	//设置api返回选项
	if apiOption != nil {
		api.SetupOption(*apiOption)
	}

	if len(conf.Addr) == 0 {
		conf.Addr = ":8080"
	}

	//手动监听服务并检测退出信号从而实现优雅退出
	srv := &http.Server{
		Addr:    conf.Addr,
		Handler: ginEngine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Info("Http listen error", zap.String("err", err.Error()))
		}
	}()

	<-errChan
	logger.GetLogger().Info("Http Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.GetLogger().Error("Http Server forced to shutdown", zap.String("err", err.Error()))
	}

	logger.GetLogger().Info("Http Server exiting")
}
