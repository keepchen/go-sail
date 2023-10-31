package sail

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/config"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/db"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/lib/nats"
	"github.com/keepchen/go-sail/v3/lib/redis"
	"github.com/keepchen/go-sail/v3/sail/httpserver"
	"go.uber.org/zap"
)

// Framework 框架配置
type Framework struct {
	appName   string
	conf      *config.Config
	apiOption *api.Option
}

// WakeupHttp 唤醒http
//
// 启动前的配置准备
//
// @param appName 应用名称
//
// @param conf 配置文件
//
// @param apiOption 统一返回配置（可选）
func WakeupHttp(appName string, conf *config.Config, apiOption *api.Option) *Framework {
	return &Framework{
		appName:   appName,
		conf:      conf,
		apiOption: apiOption,
	}
}

// Launch 启动
//
// @param registerRoutes 注册路由函数
//
// @param fn 自定义处理函数（可选），注意自定义函数是同步执行的
func (f *Framework) Launch(registerRoutes func(ginEngine *gin.Engine), fn func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.GetLogger().Error("---- Recovered ----", zap.Any("error", err))
		}
	}()

	wg := &sync.WaitGroup{}

	//:: 根据配置一次初始化组件、启动服务 ::
	//- logger
	logger.InitLoggerZap(f.conf.LoggerConf, f.appName)

	//- redis(cluster or standalone)
	if len(f.conf.RedisConf.Host) != 0 {
		redis.InitRedis(f.conf.RedisConf)
	}
	if len(f.conf.RedisClusterConf.AddrList) != 0 {
		redis.InitRedisCluster(f.conf.RedisClusterConf)
	}

	//- database
	if len(f.conf.DBConf.DriverName) != 0 {
		db.InitDB(f.conf.DBConf)
	}

	//- jwt
	if len(f.conf.JwtConf.PublicKey) != 0 {
		f.conf.JwtConf.Load()
	}

	//- nats
	if len(f.conf.NatsConf.Servers) != 0 {
		nats.Init(f.conf.NatsConf)
	}

	//- gin
	ginEngine := httpserver.InitGinEngine(f.conf.HttpServer)

	//- 注册自定义路由
	if registerRoutes != nil {
		registerRoutes(ginEngine)
	}

	//- pprof
	httpserver.EnablePProfOnDebugMode(f.conf.HttpServer, ginEngine)

	//- prometheus
	httpserver.RunPrometheusServer(f.conf.HttpServer.Prometheus)

	//- swagger
	httpserver.RunSwaggerServer(f.conf.HttpServer.Swagger, ginEngine)

	//- http server
	wg.Add(1)
	go httpserver.RunHttpServer(f.conf.HttpServer, ginEngine, f.apiOption, wg)

	printSummaryInfo(f.conf.HttpServer, ginEngine)

	//- fn
	if fn != nil {
		fn()
	}

	wg.Wait()
}
