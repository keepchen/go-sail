package routes

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/keepchen/go-sail/pkg/common/http/api"
	"github.com/keepchen/go-sail/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/pkg/app/user/config"
	"github.com/keepchen/go-sail/pkg/lib/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

const AnotherErrNone = constants.CodeType(200)

// RunServer 启动路由服务
func RunServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	var r *gin.Engine
	if config.GetGlobalConfig().Debug {
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
		r = gin.New()
	}

	//设置api返回选项
	api.SetupOption(api.Option{
		EmptyDataStruct: api.DefaultEmptyDataStructObject,
		ErrNoneCode:     AnotherErrNone,
		ErrNoneCodeMsg:  "SUCCEED",
	})

	//注册路由
	registerRoutes(ctx, r)

	//启动prometheus指标收集
	runPrometheusServerWhileEnable(ctx)

	//swagger 必须放在所有路由定义的最后
	runSwaggerServerOnDebugMode(ctx, r)

	//手动监听服务并检测退出信号从而实现优雅退出
	srv := &http.Server{
		Addr:    config.GetGlobalConfig().HttpServer.Addr,
		Handler: r,
	}

	go func() {
		log.Printf(":::::: Server listening at: %s ::::::\n", config.GetGlobalConfig().HttpServer.Addr)
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

func runSwaggerServerOnDebugMode(_ context.Context, r *gin.Engine) {
	if !config.GetGlobalConfig().HttpServer.EnableSwagger {
		//如果不是调试模式就不注册swagger路由
		return
	}

	//swagger-ui
	r.StaticFile("/swagger-assets/doc.json", "pkg/app/user/http/docs/swagger.json")
	url := ginSwagger.URL("/swagger-assets/doc.json") // The url pointing to API definition
	//access /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//redoc-ui
	r.StaticFile("/redoc/apidoc.html", "pkg/app/user/http/docs/apidoc.html")

	//favicon
	r.StaticFile("/favicon.ico", "sailboat-solid.svg")
}

// prometheus指标收集服务
//
// 当配置文件指明启用时才会启动
func runPrometheusServerWhileEnable(_ context.Context) {
	if !config.GetGlobalConfig().HttpServer.Prometheus.Enable {
		return
	}
	log.Printf("Prometheus server is ENABLE, listening address: %s\n", config.GetGlobalConfig().HttpServer.Prometheus.Addr)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		panic(http.ListenAndServe(config.GetGlobalConfig().HttpServer.Prometheus.Addr, nil))
	}()
}
