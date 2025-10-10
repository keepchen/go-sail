package go_sail

import (
	"fmt"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/sail"
	"github.com/keepchen/go-sail/v3/sail/config"
)

// ----------- 共享变量 -------------
var (
	listenerSail net.Listener
	listenerGin  net.Listener

	serviceReadySail = make(chan struct{})
	serviceReadyGin  = make(chan struct{})
)

// TestMain: 启动 go-sail 与原生 gin 两个服务
func TestMain(m *testing.M) {
	var err error
	// 启动 go-sail
	listenerSail, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get free port for go-sail: %v\n", err)
		os.Exit(1)
	}
	addrSail := normalizePort(listenerSail.Addr().String())

	go func() {
		conf := config.Config{
			LoggerConf: logger.Conf{
				Filename: "examples/logs/go-sail-benchmark.log",
			},
			HttpServer: config.HttpServerConf{
				Debug: false,
				Addr:  addrSail,
			},
		}
		register := func(r *gin.Engine) {
			r.GET("/benchmark", func(c *gin.Context) {
				sail.Response(c).Wrap(constants.ErrNone, nil).Send()
			})
		}
		after := func() {
			serviceReadySail <- struct{}{}
			close(serviceReadySail)
		}

		sail.WakeupHttp("go-sail-benchmark", &conf).
			Hook(register, nil, after).
			Launch()
	}()

	// 启动原生 Gin
	listenerGin, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get free port for gin: %v\n", err)
		os.Exit(1)
	}
	addrGin := normalizePort(listenerGin.Addr().String())

	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.GET("/benchmark", func(c *gin.Context) {
			c.JSON(http.StatusOK, dto.Base{})
		})
		server := &http.Server{
			Addr:    addrGin,
			Handler: r,
		}
		// 启动时标记就绪
		close(serviceReadyGin)
		if err := server.Serve(listenerGin); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "Gin server error: %v\n", err)
		}
	}()

	// 等待两个服务都就绪
	<-serviceReadySail
	<-serviceReadyGin

	// 跑基准
	code := m.Run()

	// 关闭 go-sail
	sail.Shutdown()

	os.Exit(code)
}

func normalizePort(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) > 1 {
		return fmt.Sprintf(":%s", parts[len(parts)-1])
	}
	return addr
}

func newBenchmarkClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
		},
		Timeout: 5 * time.Second,
	}
}

// ----------- go-sail Benchmark -----------
func BenchmarkGoSailParallel(b *testing.B) {
	serverAddr := normalizePort(listenerSail.Addr().String())
	client := newBenchmarkClient()
	url := fmt.Sprintf("http://localhost%s/benchmark", serverAddr)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(url)
			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				b.Errorf("Unexpected status: %d", resp.StatusCode)
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	})
}

// ----------- Gin Benchmark -----------
func BenchmarkGinParallel(b *testing.B) {
	serverAddr := normalizePort(listenerGin.Addr().String())
	client := newBenchmarkClient()
	url := fmt.Sprintf("http://localhost%s/benchmark", serverAddr)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, err := client.Get(url)
			if err != nil {
				b.Errorf("Request failed: %v", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				b.Errorf("Unexpected status: %d", resp.StatusCode)
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}
	})
}
