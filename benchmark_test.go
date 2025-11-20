package go_sail

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/pojo/dto"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"github.com/keepchen/go-sail/v3/sail"
	"github.com/keepchen/go-sail/v3/sail/config"
)

// ----------- 共享变量 -------------
var (
	listeningAddrSail = ":12026"
	listeningAddrGin  = ":12027"

	serviceReadySail = make(chan struct{})
	serviceReadyGin  = make(chan struct{})
)

// TestMain: 启动 go-sail 与原生 gin 两个服务
func TestMain(m *testing.M) {
	// 启动 go-sail
	go func() {
		conf := config.Config{
			LoggerConf: logger.Conf{
				Filename: "examples/logs/go-sail-benchmark.log",
			},
			HttpServer: config.HttpServerConf{
				Debug: false,
				Addr:  listeningAddrSail,
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
	go func() {
		gin.SetMode(gin.ReleaseMode)
		r := gin.New()
		r.GET("/benchmark", func(c *gin.Context) {
			c.JSON(http.StatusOK, dto.Base{})
		})
		listener, err := net.Listen("tcp", listeningAddrGin)
		if err != nil {
			panic(err)
		}
		server := &http.Server{
			Addr:    listeningAddrGin,
			Handler: r,
		}
		// 启动时标记就绪
		close(serviceReadyGin)
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			_, _ = fmt.Fprintf(os.Stderr, "Gin server error: %v\n", err)
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
	client := newBenchmarkClient()
	url := fmt.Sprintf("http://localhost%s/benchmark", listeningAddrSail)

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
			_ = resp.Body.Close()
		}
	})
}

// ----------- Gin Benchmark -----------
func BenchmarkGinParallel(b *testing.B) {
	client := newBenchmarkClient()
	url := fmt.Sprintf("http://localhost%s/benchmark", listeningAddrGin)

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
			_ = resp.Body.Close()
		}
	})
}
