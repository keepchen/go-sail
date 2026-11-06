package sail

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/sail/config"
	"github.com/keepchen/go-sail/v3/sail/httpserver"
	"go.uber.org/zap"
)

// Sailor crew ready
type Sailor interface {
	// SetupApiOption
	//
	// Set common response configuration
	SetupApiOption(opt *api.Option) Sailor
	// EnableWebsocket starts the WebSocket service
	//
	// ws WebSocket connection instance; if nil, start a connection with default settings
	//
	// handler handler function; if nil, use the default handler (only prints received messages)
	//
	// middlewares route middleware
	EnableWebsocket(ws *websocket.Conn, handler func(ws *websocket.Conn), middlewares ...gin.HandlerFunc) Sailor
	// Hook attaches related functions
	//
	// registerRoutes route registration function
	//
	// beforeFunc custom pre-hook (optional), executed before framework functions
	//
	// Note that beforeFunc runs asynchronously and components are not yet initialized; accessing components in this function will cause a nil pointer panic.
	//
	// afterFunc custom post-hook (optional), executed after framework functions
	//
	// Note that afterFunc runs asynchronously; components have been initialized according to configuration and can be used as needed.
	Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad
}

type websocketConf struct {
	enable      bool
	routePath   string
	ws          *websocket.Conn
	middlewares []gin.HandlerFunc
	handler     func(ws *websocket.Conn)
}

// Sail framework configuration
type Sail struct {
	appName   string
	conf      *config.Config
	apiOption *api.Option
	wsConf    *websocketConf
}

var _ Sailor = &Sail{}

// Launchpad launchpad
type Launchpad interface {
	// Launch launches
	Launch()
}

// Launcher launcher
type Launcher struct {
	sa                 *Sail
	registerRoutesFunc func(ginEngine *gin.Engine)
	beforeFunc         func()
	afterFunc          func()
}

var _ Launchpad = &Launcher{}

// WakeupHttp wakes up HTTP
//
// Prepare configuration before startup
//
// appName application name
//
// conf configuration file
func WakeupHttp(appName string, conf *config.Config) Sailor {
	if conf == nil {
		panic("config is nil, did you forget to initialize the config?")
	}
	return &Sail{
		appName: appName,
		conf:    conf,
	}
}

// SetupApiOption
//
// Set common response configuration
func (s *Sail) SetupApiOption(opt *api.Option) Sailor {
	s.apiOption = opt

	return s
}

// EnableWebsocket starts the WebSocket service
//
// routePath route path
//
// ws WebSocket connection instance; if nil, start the default connection instance
//
// handler handler function; if nil, start the default handler `defaultWebsocketHandlerFunc`
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

// Hook attaches related functions
//
// registerRoutes route registration function
//
// beforeFunc custom pre-hook (optional), executed before framework functions
//
// Note that beforeFunc runs asynchronously and components are not yet initialized; accessing components in this function will cause a nil pointer panic.
//
// afterFunc custom post-hook (optional), executed after framework functions
//
// Note that afterFunc runs asynchronously; components have been initialized according to configuration and can be used as needed.
func (s *Sail) Hook(registerRoutes func(ginEngine *gin.Engine), beforeFunc, afterFunc func()) Launchpad {
	return &Launcher{
		sa:                 s,
		registerRoutesFunc: registerRoutes,
		beforeFunc:         beforeFunc,
		afterFunc:          afterFunc,
	}
}

// Launch launches
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

	//- before, invoke the custom pre-hook
	if l.beforeFunc != nil {
		go l.beforeFunc()
	}

	//:: Initialize components and start services in order according to configuration ::
	componentsStartup(l.sa.appName, l.sa.conf)

	//- Gin
	ginEngine := httpserver.InitGinEngine(l.sa.conf.HttpServer)

	//- Register custom routes
	if l.registerRoutesFunc != nil {
		l.registerRoutesFunc(ginEngine)
	}

	//- Register WebSocket
	if l.sa.wsConf != nil && l.sa.wsConf.enable {
		if len(l.sa.wsConf.middlewares) > 0 {
			ginEngine.Use(l.sa.wsConf.middlewares...).GET(l.sa.wsConf.routePath,
				httpserver.WrapWebsocketHandler(l.sa.wsConf.ws, l.sa.wsConf.handler))
		} else {
			ginEngine.GET(l.sa.wsConf.routePath, httpserver.WrapWebsocketHandler(l.sa.wsConf.ws, l.sa.wsConf.handler))
		}
	}

	//- Pprof
	httpserver.EnablePProfOnDebugMode(l.sa.conf.HttpServer, ginEngine)

	//- Prometheus
	httpserver.RunPrometheusServerWhenEnable(l.sa.conf.HttpServer.Prometheus)

	//- Swagger
	httpserver.RunSwaggerServerWhenEnable(l.sa.conf.HttpServer.Swagger, ginEngine)

	//- http server
	wg.Add(1)
	go httpserver.RunHttpServer(l.sa.conf.HttpServer, ginEngine, l.sa.apiOption, wg)

	printSummaryInfo(*l.sa.conf, ginEngine)

	//set config
	config.Set(l.sa.conf)

	//- after, invoke the custom post-hook
	if l.afterFunc != nil {
		go l.afterFunc()
	}

	wg.Wait()

	//Waiting for the caller before component shutdown is enabled; wait for the caller's exit signal
	if atomic.LoadInt32(&componentsWaitForCallerState) == 1 {
		<-waitForMeChan
	}

	//Shut down the corresponding components
	componentsShutdown(l.sa.conf)
}

// Shutdown shuts down services
//
// Used to stop services manually
func Shutdown() {
	go httpserver.NotifyExit(syscall.SIGINT)
}

var (
	componentsWaitForCallerState int32
	waitForMeChan                = make(chan struct{})
	closeWaitChanState           int32
)

// ComponentsBeforeShutdownWaitForCaller waits for the caller before shutting down components
//
// Default framework behavior: stop the HTTP service and shut down started components upon receiving an exit signal
//
// After this method is called, shutting down started components waits for the caller's instruction;
// components are shut down only after that instruction is received.
//
// Callers can invoke ComponentsShouldShutdownCauseImReady to signal component shutdown.
//
// # Usage example:
//
//		  beforeFunc := func() {
//	       sail.ComponentsBeforeShutdownWaitForCaller()
//		  }
//
//			 sail.WakeupHttp("go-sail", conf).
//					Hook(nil, beforeFunc, nil).
//					Launch()
func ComponentsBeforeShutdownWaitForCaller() {
	atomic.CompareAndSwapInt32(&componentsWaitForCallerState, 0, 1)
}

// ComponentsShouldShutdownCauseImReady signals readiness for component shutdown
//
// Use this method together with ComponentsBeforeShutdownWaitForCaller
func ComponentsShouldShutdownCauseImReady() {
	if atomic.CompareAndSwapInt32(&closeWaitChanState, 0, 1) {
		close(waitForMeChan)
	}
}
