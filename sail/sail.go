package sail

import (
	"sync"

	"github.com/keepchen/go-sail/v3/lib/kafka"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/gin-gonic/gin"
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
// @param before 前置自定义处理函数（可选），在框架函数之前执行，注意自定义函数是同步执行的
//
// @param after 后置自定义处理函数（可选），在框架函数之后执行，注意自定义函数是同步执行的
func (f *Framework) Launch(registerRoutes func(ginEngine *gin.Engine), before, after func()) {
	defer func() {
		if err := recover(); err != nil {
			logger.GetLogger().Error("---- Recovered ----", zap.Any("error", err))
		}
	}()

	wg := &sync.WaitGroup{}

	//:: 根据配置一次初始化组件、启动服务 ::
	//- before
	if before != nil {
		before()
	}

	//- logger
	logger.InitLoggerZap(f.conf.LoggerConf, f.appName)

	//- redis(standalone)
	if len(f.conf.RedisConf.Host) != 0 {
		redis.InitRedis(f.conf.RedisConf)
	}

	//- redis(cluster)
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

	//- kafka
	if len(f.conf.KafkaConf.Conf.AddrList) != 0 {
		kafka.Init(f.conf.KafkaConf.Conf, f.conf.KafkaConf.Topic, f.conf.KafkaConf.GroupID)
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

	//- after
	if after != nil {
		after()
	}

	wg.Wait()
}
