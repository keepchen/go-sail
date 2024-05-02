package sail

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"

	"github.com/keepchen/go-sail/v3/lib/logger"
	"go.uber.org/zap"

	"github.com/gorilla/websocket"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/sail/httpserver"
)

// Sailor 船员就位
type Sailor interface {
	// SetupApiOption
	//
	// 设置统一返回配置
	SetupApiOption(opt *api.Option) Sailor
	// EnableWebsocket 启动websocket服务
	//
	// ws websocket连接实例，若为空，则启动默认配置连接
	//
	// handler 处理函数，若为空，则启用默认处理函数（仅打印接收到的message信息）
	//
	// middlewares 路由中间件
	EnableWebsocket(ws *websocket.Conn, handler func(ws *websocket.Conn), middlewares ...gin.HandlerFunc) Sailor
	// Hook 挂载相关方法
	//
	// registerRoutes 注册路由函数
	//
	// beforeFunc 前置自定义处理函数（可选），在框架函数之前执行
	//
	//注意beforeFunc是异步执行的，另外此时组件尚未初始化，因此在此函数内调用组件将会出现空指针异常（panic）。
	//
	// afterFunc 后置自定义处理函数（可选），在框架函数之后执行
	//
	//注意afterFunc是异步执行的，另外此时组件已经按配置初始化完成，可以按需调用。
	Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad
}

type websocketConf struct {
	enable      bool
	routePath   string
	ws          *websocket.Conn
	middlewares []gin.HandlerFunc
	handler     func(ws *websocket.Conn)
}

// Sail 框架配置
type Sail struct {
	appName   string
	conf      *config.Config
	apiOption *api.Option
	wsConf    *websocketConf
}

var _ Sailor = &Sail{}

// Launchpad 启动台
type Launchpad interface {
	// Launch 启动
	Launch()
}

// Launcher 启动器
type Launcher struct {
	sa                 *Sail
	registerRoutesFunc func(ginEngine *gin.Engine)
	beforeFunc         func()
	afterFunc          func()
}

var _ Launchpad = &Launcher{}

// WakeupHttp 唤醒http
//
// 启动前的配置准备
//
// appName 应用名称
//
// conf 配置文件
func WakeupHttp(appName string, conf *config.Config) Sailor {
	return &Sail{
		appName: appName,
		conf:    conf,
	}
}

// SetupApiOption
//
// 设置统一返回配置
func (s *Sail) SetupApiOption(opt *api.Option) Sailor {
	s.apiOption = opt

	return s
}

// EnableWebsocket 启动websocket服务
//
// routePath 路由地址
//
// ws websocket连接实例，若为空，则启动默认的连接实例
//
// handler 处理函数，若为空，则启动`defaultWebsocketHandlerFunc`默认处理函数
func (s *Sail) EnableWebsocket(ws *websocket.Conn, handler func(ws *websocket.Conn), middlewares ...gin.HandlerFunc) Sailor {
	s.wsConf = &websocketConf{
		enable:      true,
		routePath:   s.conf.HttpServer.WebSocketRoutePath,
		ws:          ws,
		middlewares: middlewares,
		handler:     handler,
	}

	return s
}

// Hook 挂载相关方法
//
// registerRoutes 注册路由函数
//
// beforeFunc 前置自定义处理函数（可选），在框架函数之前执行
//
// 注意beforeFunc是异步执行的，另外此时组件尚未初始化，因此在此函数内调用组件将会出现空指针异常（panic）。
//
// afterFunc 后置自定义处理函数（可选），在框架函数之后执行
//
// 注意afterFunc是异步执行的，另外此时组件已经按配置初始化完成，可以按需调用。
func (s *Sail) Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad {
	return &Launcher{
		sa:                 s,
		registerRoutesFunc: registerRoutes,
		beforeFunc:         beforeFunc,
		afterFunc:          afterFunc,
	}
}

// Launch 启动
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

	//- before，自定义前置函数调用
	if l.beforeFunc != nil {
		go l.beforeFunc()
	}

	//:: 根据配置依次初始化组件、启动服务 ::
	componentsStartup(l.sa.appName, l.sa.conf)

	//- gin
	ginEngine := httpserver.InitGinEngine(l.sa.conf.HttpServer)

	//- 注册自定义路由
	if l.registerRoutesFunc != nil {
		l.registerRoutesFunc(ginEngine)
	}

	//- 注册websocket
	if l.sa.wsConf != nil && l.sa.wsConf.enable {
		if len(l.sa.wsConf.middlewares) > 0 {
			ginEngine.Use(l.sa.wsConf.middlewares...).GET(l.sa.wsConf.routePath, httpserver.WrapWebsocketHandler(l.sa.wsConf.ws, l.sa.wsConf.handler))
		} else {
			ginEngine.GET(l.sa.wsConf.routePath, httpserver.WrapWebsocketHandler(l.sa.wsConf.ws, l.sa.wsConf.handler))
		}
	}

	//- pprof
	httpserver.EnablePProfOnDebugMode(l.sa.conf.HttpServer, ginEngine)

	//- prometheus
	httpserver.RunPrometheusServerWhenEnable(l.sa.conf.HttpServer.Prometheus)

	//- swagger
	httpserver.RunSwaggerServerWhenEnable(l.sa.conf.HttpServer.Swagger, ginEngine)

	//- http server
	wg.Add(1)
	go httpserver.RunHttpServer(l.sa.conf.HttpServer, ginEngine, l.sa.apiOption, wg)

	printSummaryInfo(l.sa.conf.HttpServer, ginEngine)

	//- after,自定义后置函数调用
	if l.afterFunc != nil {
		go l.afterFunc()
	}

	wg.Wait()

	//关闭响应的组件
	componentsShutdown(l.sa.conf)
}
