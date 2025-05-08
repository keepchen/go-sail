package sail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"syscall"
	"testing"
	"time"

	"go.uber.org/zap"

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

const testServerAddr = ":12345"

func TestLogTrace(t *testing.T) {
	conf := &config.Config{
		LoggerConf: logger.Conf{
			Level:    "debug",
			Filename: "../examples/logs/testcase_logtrace.log",
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
		Hook(registerRoutes4LogTrace, nil, afterFunc).
		Launch()
}

func registerRoutes4LogTrace(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		fmt.Println(LogTrace(c).RequestID())
		fmt.Println(LogTrace(c).SpanID())
		LogTrace(c).GetLogger().Debug("registerRoutes4LogTrace/ping/GetLogger", zap.String("value", "ok"))
		LogTrace(c).Debug("registerRoutes4LogTrace/ping/Debug", zap.String("value", "ok"))
		LogTrace(c).Info("registerRoutes4LogTrace/ping/Info", zap.String("value", "ok"))
		LogTrace(c).Warn("registerRoutes4LogTrace/ping/Warn", zap.String("value", "ok"))
		LogTrace(c).Error("registerRoutes4LogTrace/ping/Error", zap.String("value", "ok"))
		c.String(http.StatusOK, "%s", "pong")
	}).
		GET("/hello", func(c *gin.Context) {
			Response(c).Wrap(constants.CodeType(200), "world").Send()
		})
}
