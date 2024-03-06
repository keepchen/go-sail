package sail

import (
	"fmt"
	"runtime"
	"runtime/debug"
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
	// SetupApiOption
	//
	// 设置统一返回配置
	SetupApiOption(opt *api.Option) Sailor
	// Hook 挂载相关方法
	//
	// @param registerRoutes 注册路由函数
	//
	// @param beforeFunc 前置自定义处理函数（可选），在框架函数之前执行，注意自定义函数是同步执行的
	//
	// @param afterFunc 后置自定义处理函数（可选），在框架函数之后执行，注意自定义函数是同步执行的
	Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad
	// Launch 启动
	//
	// @param registerRoutes 注册路由函数
	//
	// # Note:
	//
	// 未设置前置自动函数、未设置后置自定义函数
	Launch(registerRoutes func(ginEngine *gin.Engine))
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
	// Launch 启动
	//
	// # Note:
	//
	// 已注册路由、已设置前置自动函数、已设置后置自定义函数
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
			pc := make([]uintptr, 10)
			n := runtime.Callers(3, pc)
			frames := runtime.CallersFrames(pc[:n])
			frame, _ := frames.Next()
			fmt.Printf("[GO-SAIL] Try to recover but failed\nReason: %v\nCaller: %s:%d -> %s\nStack:\n",
				err, frame.File, frame.Line, frame.Function)
			debug.PrintStack()
			logger.GetLogger().Error("---- Recovered ----", zap.Any("error", err))
		}
	}()

	wg := &sync.WaitGroup{}

	//- logger
	logger.Init(f.conf.LoggerConf, f.appName)

	//- redis(standalone)
	if f.conf.RedisConf.Enable {
		redis.InitRedis(f.conf.RedisConf)
	}

	//- redis(cluster)
	if f.conf.RedisClusterConf.Enable {
		redis.InitRedisCluster(f.conf.RedisClusterConf)
	}

	//- database
	if f.conf.DBConf.Enable {
		db.Init(f.conf.DBConf)
	}

	//- jwt
	if f.conf.JwtConf.Enable {
		f.conf.JwtConf.Load()
	}

	//- nats
	if f.conf.NatsConf.Enable {
		nats.Init(f.conf.NatsConf)
	}

	//- kafka
	if f.conf.KafkaConf.Conf.Enable {
		kafka.Init(f.conf.KafkaConf.Conf, f.conf.KafkaConf.Topic, f.conf.KafkaConf.GroupID)
	}

	//- etcd
	if f.conf.EtcdConf.Enable {
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
			pc := make([]uintptr, 10)
			n := runtime.Callers(3, pc)
			frames := runtime.CallersFrames(pc[:n])
			frame, _ := frames.Next()
			fmt.Printf("[GO-SAIL] Try to recover but failed\nReason: %v\nCaller: %s:%d -> %s\nStack:\n",
				err, frame.File, frame.Line, frame.Function)
			debug.PrintStack()
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
	if l.fw.conf.RedisConf.Enable {
		redis.InitRedis(l.fw.conf.RedisConf)
	}

	//- redis(cluster)
	if l.fw.conf.RedisClusterConf.Enable {
		redis.InitRedisCluster(l.fw.conf.RedisClusterConf)
	}

	//- database
	if l.fw.conf.DBConf.Enable {
		db.Init(l.fw.conf.DBConf)
	}

	//- jwt
	if l.fw.conf.JwtConf.Enable {
		l.fw.conf.JwtConf.Load()
	}

	//- nats
	if l.fw.conf.NatsConf.Enable {
		nats.Init(l.fw.conf.NatsConf)
	}

	//- kafka
	if l.fw.conf.KafkaConf.Conf.Enable {
		kafka.Init(l.fw.conf.KafkaConf.Conf, l.fw.conf.KafkaConf.Topic, l.fw.conf.KafkaConf.GroupID)
	}

	//- etcd
	if l.fw.conf.EtcdConf.Enable {
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
