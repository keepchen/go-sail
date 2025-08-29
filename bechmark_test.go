package go_sail

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/sail"
	"github.com/keepchen/go-sail/v3/sail/config"
)

// listener 用于避免端口冲突
var (
	listener     net.Listener
	serviceReady = make(chan struct{})
)

// TestMain 用于初始化和清理测试环境
func TestMain(m *testing.M) {
	// 获取可用端口 listener
	var err error
	listener, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get free port: %v\n", err)
		os.Exit(1)
	}

	serverAddr := listener.Addr().String()

	addrSplit := strings.Split(serverAddr, ":")
	if len(addrSplit) > 1 {
		serverAddr = fmt.Sprintf(":%s", addrSplit[1])
	}

	// 启动 go-sail 服务
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		conf := config.Config{
			LoggerConf: logger.Conf{
				Filename: "examples/logs/go-sail-benchmark.log",
			},
			HttpServer: config.HttpServerConf{
				Debug: false,
				Addr:  serverAddr,
			},
		}

		registerRoutes := func(r *gin.Engine) {
			r.GET("/benchmark", func(context *gin.Context) {
				sail.Response(context).Wrap(constants.ErrNone, nil).Send()
			})
		}

		afterFunc := func() {
			serviceReady <- struct{}{}
			close(serviceReady)
		}

		sail.WakeupHttp("go-sail-benchmark", &conf).
			Hook(registerRoutes, nil, afterFunc).
			Launch()
	}()

	// 等待服务器完全启动
	<-serviceReady

	// 运行 benchmark
	code := m.Run()

	// 清理：关闭服务器
	sail.Shutdown()
	wg.Wait()

	os.Exit(code)
}

// BenchmarkGoSailParallel 并发 HTTP 请求基准测试
func BenchmarkGoSailParallel(b *testing.B) {
	serverAddr := listener.Addr().String()
	addrSplit := strings.Split(serverAddr, ":")
	if len(addrSplit) > 1 {
		serverAddr = fmt.Sprintf(":%s", addrSplit[1])
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
		},
		Timeout: 5 * time.Second,
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(fmt.Sprintf("http://localhost%s/benchmark", serverAddr))
			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				b.Errorf("Unexpected status code: %d", resp.StatusCode)
			}

			// 读取响应体但不保存，保证连接复用
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	})
}
