package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusExporter 导出Prometheus指标
func PrometheusExporter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("entryUnixMilli", time.Now().UnixMilli())

		w := responseBodyWriter{ctx: c, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	ctx *gin.Context
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	entryUnixMilli := r.ctx.MustGet("entryUnixMilli").(int64)
	elapsed := time.Now().UnixMilli() - entryUnixMilli
	//记录接口延时
	metricsSummaryVecLatency.WithLabelValues(r.ctx.FullPath()).Observe(float64(elapsed))
	//记录接口返回状态码
	metricsSummaryVecHttpStatus.WithLabelValues(r.ctx.FullPath(), fmt.Sprintf("%d", r.ctx.Writer.Status())).Observe(float64(1))

	return r.ResponseWriter.Write(b)
}

var (
	once                        sync.Once
	metricsSummaryVecLatency    *prometheus.SummaryVec
	metricsSummaryVecHttpStatus *prometheus.SummaryVec
)

func init() {
	once.Do(func() {
		svl := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "api_durations",
			Help: "[api] http latency distributions (milliseconds)",
		},
			[]string{"path"},
		)

		svh := prometheus.NewSummaryVec(prometheus.SummaryOpts{
			Name: "api_http_status",
			Help: "[api] http status code",
		},
			[]string{"path", "status"},
		)

		prometheus.MustRegister(svl, svh)

		metricsSummaryVecLatency = svl
		metricsSummaryVecHttpStatus = svh
	})
}
