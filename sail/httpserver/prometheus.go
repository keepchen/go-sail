package httpserver

import (
	"net/http"

	"github.com/keepchen/go-sail/v3/sail/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RunPrometheusServerWhenEnable 启动prometheus指标收集服务
//
// 当配置文件指明启用时才会启动
func RunPrometheusServerWhenEnable(conf config.PrometheusConf) {
	if !conf.Enable {
		return
	}
	var path = conf.AccessPath
	if len(path) == 0 {
		path = "/metrics"
	}
	go func() {
		http.Handle(path, promhttp.Handler())
		panic(http.ListenAndServe(conf.Addr, nil))
	}()
}
