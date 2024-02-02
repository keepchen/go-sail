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

// InitGinEngine 初始化gin引擎
//
// # Note:
//
// 该方法会默认使用 RequestEntry 中间件
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

	r.Use(middleware.RequestEntry())

	return r
}

// RunHttpServer 启动http服务
func RunHttpServer(conf config.HttpServerConf, ginEngine *gin.Engine, apiOption *api.Option, wg *sync.WaitGroup) {
	defer wg.Done()

	//监听退出信号
	errChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%v", <-c)
		cancel()
	}()

	//设置api返回选项
	if apiOption != nil {
		api.SetupOption(*apiOption)
	}

	//手动监听服务并检测退出信号从而实现优雅退出
	srv := &http.Server{
		Addr:    conf.Addr,
		Handler: ginEngine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.GetLogger().Info("listen error", zap.Errors("error", []error{err}))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.GetLogger().Info("Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.GetLogger().Error("Server forced to shutdown", zap.Errors("errors", []error{err}))
	}

	logger.GetLogger().Info("Server exiting")
}
