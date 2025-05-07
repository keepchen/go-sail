package httpserver

import (
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/logger"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/sail/config"
)

var loggerConf = logger.Conf{
	Filename: "../../examples/logs/httpserver_tester.log",
}

func TestInitGinEngine(t *testing.T) {
	t.Run("InitGinEngine", func(t *testing.T) {
		conf := config.HttpServerConf{}
		InitGinEngine(conf)
	})

	t.Run("InitGinEngine-Debug", func(t *testing.T) {
		conf := config.HttpServerConf{
			Debug: true,
		}
		InitGinEngine(conf)
	})

	t.Run("InitGinEngine-TrustedProxies", func(t *testing.T) {
		conf := config.HttpServerConf{
			TrustedProxies: []string{"10.0.0.0/16"},
		}
		InitGinEngine(conf)
	})

	t.Run("InitGinEngine-PrometheusEnable", func(t *testing.T) {
		conf := config.HttpServerConf{
			Prometheus: config.PrometheusConf{
				Enable: true,
			},
		}
		InitGinEngine(conf)
	})
}

func TestRunHttpServer(t *testing.T) {
	t.Run("RunHttpServer", func(t *testing.T) {
		go func() {
			time.Sleep(5 * time.Second)
			NotifyExit(os.Interrupt)
		}()

		logger.Init(loggerConf, "go-sail")

		conf := config.HttpServerConf{}
		r := InitGinEngine(conf)
		wg := &sync.WaitGroup{}
		wg.Add(1)
		RunHttpServer(conf, r, &api.Option{}, wg)
	})
}
