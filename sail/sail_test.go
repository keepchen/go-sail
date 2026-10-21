package sail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/http/pojo/dto"
	"github.com/keepchen/go-sail/v3/lib/logger"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/sail/httpserver"
	"github.com/keepchen/go-sail/v3/utils"
	"github.com/stretchr/testify/assert"

	"github.com/gin-gonic/gin"

	"github.com/keepchen/go-sail/v3/sail/config"
)

func TestWakeHttp(t *testing.T) {
	t.Run("wakeHttpPanic", func(t *testing.T) {
		assert.Panics(t, func() {
			WakeupHttp("test-wakeup-panic", nil)
		})
	})
	t.Run("wakeHttp", func(t *testing.T) {
		t.Log(WakeupHttp("test-wakeup", &config.Config{}))
	})
}

func TestLaunch(t *testing.T) {
	conf := &config.Config{
		LoggerConf: logger.Conf{
			Level:    "debug",
			Filename: "../examples/logs/testcase_wakeupHttp.log",
		},
		HttpServer: config.HttpServerConf{
			Addr: testServerAddr,
		},
	}
	apiOption := api.DefaultSetupOption()
	apiOption.ErrNoneCode = constants.CodeType(200)

	beforeFunc := func() {}

	afterFunc := func() {
		fmt.Println("service is ready")
	}

	//模拟请求
	requestFunc := func() {
		timeout := time.Second * 5
		//ping
		t.Run("requestPing", func(t *testing.T) {
			url := fmt.Sprintf("http://localhost%s/ping", conf.HttpServer.Addr)
			resp, statusCode, err := utils.HttpClient().SendRequest("GET", url, nil, nil, timeout)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, statusCode)
			assert.Equal(t, "pong", string(resp))
		})

		//hello
		t.Run("requestHello", func(t *testing.T) {
			url := fmt.Sprintf("http://localhost%s/hello", conf.HttpServer.Addr)
			resp, statusCode, err := utils.HttpClient().SendRequest("GET", url, nil, nil, timeout)
			var ack dto.Base
			_ = json.Unmarshal(resp, &ack)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, statusCode)
			assert.Equal(t, "world", ack.Data)
			assert.Equal(t, apiOption.ErrNoneCode.Int(), ack.Code)
		})
	}

	//手动退出
	go func() {
		timer := time.NewTimer(time.Second * 5)
		for range timer.C {
			requestFunc()
			httpserver.NotifyExit(syscall.SIGTERM)
		}
	}()

	WakeupHttp("go-sail-tester", conf).
		SetupApiOption(apiOption).
		Hook(registerRoutes, beforeFunc, afterFunc).
		Launch()
}

func registerRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "pong")
	}).
		GET("/hello", func(c *gin.Context) {
			Response(c).Wrap(constants.CodeType(200), "world").Send()
		})
}

func TestEnableWebsocket(t *testing.T) {
	t.Run("EnableWebsocket", func(t *testing.T) {
		conf := &config.Config{
			LoggerConf: logger.Conf{
				Level:    "debug",
				Filename: "../examples/logs/testcase_wakeupHttp.log",
			},
			HttpServer: config.HttpServerConf{
				Addr: testServerAddr,
			},
		}
		apiOption := api.DefaultSetupOption()
		apiOption.ErrNoneCode = constants.CodeType(200)

		afterFunc := func() {
			fmt.Println("service is ready")
		}
		t.Log(WakeupHttp("go-sail-tester", conf).
			SetupApiOption(apiOption).
			EnableWebsocket(nil, nil).
			Hook(registerRoutes, nil, afterFunc))
	})
}

func TestShutdown(t *testing.T) {
	t.Run("Shutdown", func(t *testing.T) {
		Shutdown()
	})
}

func TestComponentsBeforeShutdownWaitForCaller(t *testing.T) {
	resetComponentsWaitForCallerState(t)

	ComponentsBeforeShutdownWaitForCaller()

	assert.Equal(t, int32(1), atomic.LoadInt32(&componentsWaitForCallerState))
	select {
	case <-waitForMeChan:
		t.Fatal("wait channel should remain blocked before caller is ready")
	default:
	}

	// Repeated calls should be idempotent.
	ComponentsBeforeShutdownWaitForCaller()
	assert.Equal(t, int32(1), atomic.LoadInt32(&componentsWaitForCallerState))
}

func TestComponentsShouldShutdownCauseImReady(t *testing.T) {
	t.Run("readyAfterWaitEnabled", func(t *testing.T) {
		resetComponentsWaitForCallerState(t)
		ComponentsBeforeShutdownWaitForCaller()

		const callers = 16
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(callers)
		for range callers {
			go func() {
				defer wg.Done()
				<-start
				ComponentsShouldShutdownCauseImReady()
			}()
		}
		close(start)
		wg.Wait()

		assert.Equal(t, int32(1), atomic.LoadInt32(&closeWaitChanState))
		assertWaitForMeChanClosed(t)
	})

	t.Run("readyBeforeWaitEnabled", func(t *testing.T) {
		resetComponentsWaitForCallerState(t)

		ComponentsShouldShutdownCauseImReady()
		ComponentsBeforeShutdownWaitForCaller()

		assert.Equal(t, int32(1), atomic.LoadInt32(&componentsWaitForCallerState))
		assert.Equal(t, int32(1), atomic.LoadInt32(&closeWaitChanState))
		assertWaitForMeChanClosed(t)
	})
}

func resetComponentsWaitForCallerState(t *testing.T) {
	t.Helper()
	t.Cleanup(func() {
		atomic.StoreInt32(&componentsWaitForCallerState, 0)
		atomic.StoreInt32(&closeWaitChanState, 0)
		waitForMeChan = make(chan struct{})
	})
	atomic.StoreInt32(&componentsWaitForCallerState, 0)
	atomic.StoreInt32(&closeWaitChanState, 0)
	waitForMeChan = make(chan struct{})
}

func assertWaitForMeChanClosed(t *testing.T) {
	t.Helper()
	select {
	case _, ok := <-waitForMeChan:
		assert.False(t, ok, "wait channel should be closed")
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for caller-ready signal")
	}
}
