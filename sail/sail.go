package sail

import (
	"sync"

	"github.com/keepchen/go-sail/v3/lib/etcd"

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

// Sailor 船员就位
type Sailor interface {
	SetupApiOption(opt *api.Option) Sailor
	Launch(registerRoutes func(ginEngine *gin.Engine))
	Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad
}

// Framework 框架配置
type Framework struct {
	appName   string
	conf      *config.Config
	apiOption *api.Option
}

var _ Sailor = &Framework{}

// Launchpad 启动台
type Launchpad interface {
	Launch()
}

// Launcher 启动器
type Launcher struct {
	fw                 *Framework
	registerRoutesFunc func(ginEngine *gin.Engine)
	beforeFunc         func()
	afterFunc          func()
}

var _ Launchpad = &Launcher{}

// WakeupHttp 唤醒http
//
// 启动前的配置准备
//
// @param appName 应用名称
//
// @param conf 配置文件
//
// @param apiOption 统一返回配置（可选）
func WakeupHttp(appName string, conf *config.Config) Sailor {
	return &Framework{
		appName: appName,
		conf:    conf,
	}
}

// SetupApiOption
//
// 设置统一返回配置
func (f *Framework) SetupApiOption(opt *api.Option) Sailor {
	f.apiOption = opt

	return f
}

// Launch 启动
//
// @param registerRoutes 注册路由函数
//
// # Note:
//
// 未设置前置自动函数、未设置后置自定义函数
func (f *Framework) Launch(registerRoutes func(ginEngine *gin.Engine)) {
	defer func() {
		if err := recover(); err != nil {
			logger.GetLogger().Error("---- Recovered ----", zap.Any("error", err))
		}
	}()

	wg := &sync.WaitGroup{}

	//- logger
	logger.Init(f.conf.LoggerConf, f.appName)

	//- redis(standalone)
	if len(f.conf.RedisConf.Host) > 0 {
		redis.InitRedis(f.conf.RedisConf)
	}

	//- redis(cluster)
	if len(f.conf.RedisClusterConf.AddrList) > 0 {
		redis.InitRedisCluster(f.conf.RedisClusterConf)
	}

	//- database
	if len(f.conf.DBConf.DriverName) > 0 {
		db.Init(f.conf.DBConf)
	}

	//- jwt
	if len(f.conf.JwtConf.PublicKey) > 0 {
		f.conf.JwtConf.Load()
	}

	//- nats
	if len(f.conf.NatsConf.Servers) > 0 {
		nats.Init(f.conf.NatsConf)
	}

	//- kafka
	if len(f.conf.KafkaConf.Conf.AddrList) > 0 {
		kafka.Init(f.conf.KafkaConf.Conf, f.conf.KafkaConf.Topic, f.conf.KafkaConf.GroupID)
	}

	//- etcd
	if len(f.conf.EtcdConf.Endpoints) > 0 {
		etcd.Init(f.conf.EtcdConf)
	}

	//- gin
	ginEngine := httpserver.InitGinEngine(f.conf.HttpServer)
	if registerRoutes != nil {
		registerRoutes(ginEngine)
	}

	//- pprof
	httpserver.EnablePProfOnDebugMode(f.conf.HttpServer, ginEngine)

	//- prometheus
	httpserver.RunPrometheusServerWhenEnable(f.conf.HttpServer.Prometheus)

	//- swagger
	httpserver.RunSwaggerServerWhenEnable(f.conf.HttpServer.Swagger, ginEngine)

	//- http server
	wg.Add(1)
	go httpserver.RunHttpServer(f.conf.HttpServer, ginEngine, f.apiOption, wg)

	printSummaryInfo(f.conf.HttpServer, ginEngine)

	wg.Wait()
}

// Hook 挂载相关方法
//
// @param registerRoutes 注册路由函数
//
// @param beforeFunc 前置自定义处理函数（可选），在框架函数之前执行，注意自定义函数是同步执行的
//
// @param afterFunc 后置自定义处理函数（可选），在框架函数之后执行，注意自定义函数是同步执行的
func (f *Framework) Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad {
	return &Launcher{
		fw:                 f,
		registerRoutesFunc: registerRoutes,
		beforeFunc:         beforeFunc,
		afterFunc:          afterFunc,
	}
}

// Launch 启动
//
// # Note:
//
// 已注册路由、已设置前置自动函数、已设置后置自定义函数
func (l *Launcher) Launch() {
	defer func() {
		if err := recover(); err != nil {
			logger.GetLogger().Error("---- Recovered ----", zap.Any("error", err))
		}
	}()

	wg := &sync.WaitGroup{}

	//:: 根据配置依次初始化组件、启动服务 ::
	//- before，自定义前置函数调用
	if l.beforeFunc != nil {
		l.beforeFunc()
	}

	//- logger
	logger.Init(l.fw.conf.LoggerConf, l.fw.appName)

	//- redis(standalone)
	if len(l.fw.conf.RedisConf.Host) > 0 {
		redis.InitRedis(l.fw.conf.RedisConf)
	}

	//- redis(cluster)
	if len(l.fw.conf.RedisClusterConf.AddrList) > 0 {
		redis.InitRedisCluster(l.fw.conf.RedisClusterConf)
	}

	//- database
	if len(l.fw.conf.DBConf.DriverName) > 0 {
		db.Init(l.fw.conf.DBConf)
	}

	//- jwt
	if len(l.fw.conf.JwtConf.PublicKey) > 0 {
		l.fw.conf.JwtConf.Load()
	}

	//- nats
	if len(l.fw.conf.NatsConf.Servers) > 0 {
		nats.Init(l.fw.conf.NatsConf)
	}

	//- kafka
	if len(l.fw.conf.KafkaConf.Conf.AddrList) > 0 {
		kafka.Init(l.fw.conf.KafkaConf.Conf, l.fw.conf.KafkaConf.Topic, l.fw.conf.KafkaConf.GroupID)
	}

	//- etcd
	if len(l.fw.conf.EtcdConf.Endpoints) > 0 {
		etcd.Init(l.fw.conf.EtcdConf)
	}

	//- gin
	ginEngine := httpserver.InitGinEngine(l.fw.conf.HttpServer)

	//- 注册自定义路由
	if l.registerRoutesFunc != nil {
		l.registerRoutesFunc(ginEngine)
	}

	//- pprof
	httpserver.EnablePProfOnDebugMode(l.fw.conf.HttpServer, ginEngine)

	//- prometheus
	httpserver.RunPrometheusServerWhenEnable(l.fw.conf.HttpServer.Prometheus)

	//- swagger
	httpserver.RunSwaggerServerWhenEnable(l.fw.conf.HttpServer.Swagger, ginEngine)

	//- http server
	wg.Add(1)
	go httpserver.RunHttpServer(l.fw.conf.HttpServer, ginEngine, l.fw.apiOption, wg)

	printSummaryInfo(l.fw.conf.HttpServer, ginEngine)

	//- after,自定义后置函数调用
	if l.afterFunc != nil {
		l.afterFunc()
	}

	wg.Wait()
}
