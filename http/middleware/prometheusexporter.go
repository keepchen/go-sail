package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusExporter exports Prometheus metrics
func PrometheusExporter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("entryUnixMilli", time.Now().UnixMilli())

		w := &responseBodyWriter{ctx: c, ResponseWriter: c.Writer, written: false}
		c.Writer = w

		c.Next()
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	ctx     *gin.Context
	written bool
}

func (r *responseBodyWriter) Write(b []byte) (int, error) {
	//In order to ensure that the call chain of a single request is only recorded once
	//issue such as: https://github.com/gin-contrib/gzip/issues/47
	if r.written {
		return r.ResponseWriter.Write(b)
	}

	entryUnixMilli := r.ctx.MustGet("entryUnixMilli").(int64)
	elapsed := time.Now().UnixMilli() - entryUnixMilli
	//Record API latency
	metricsSummaryVecHttpLatency.WithLabelValues(r.ctx.FullPath()).Observe(float64(elapsed))
	//Record the API response status code
	metricsSummaryVecHttpStatus.WithLabelValues(r.ctx.FullPath(), fmt.Sprintf("%d", r.ctx.Writer.Status())).Observe(float64(1))

	r.written = true

	return r.ResponseWriter.Write(b)
}

var once sync.Once

// HTTP API status and latency
var (
	metricsSummaryVecHttpLatency          *prometheus.SummaryVec
	metricsSummaryVecHttpStatus           *prometheus.SummaryVec
	metricsSummaryGaugeVecCPUUsage        *prometheus.GaugeVec
	metricsSummaryGaugeVecMemoryUsage     *prometheus.GaugeVec
	metricsSummaryGaugeVecDiskUsage       *prometheus.GaugeVec
	metricsSummaryGaugeVecNetworkTransfer *prometheus.GaugeVec
)

var (
	diskPath       = "/"
	sampleInterval = time.Minute
)

// SetDiskPath sets the disk monitoring path
func SetDiskPath(path string) {
	if len(path) != 0 {
		diskPath = path
	}
}

// SetSampleInterval sets the sampling interval (frequency)
func SetSampleInterval(interval string) {
	if len(interval) != 0 {
		sampleInterval = time.Minute
		return
	}
	td, err := time.ParseDuration(interval)
	//Use the default of 1 minute if less than 1ms
	if err != nil || td.Milliseconds() == 0 {
		td = time.Minute
	}
	sampleInterval = td
}

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

		metricsSummaryVecHttpLatency = svl
		metricsSummaryVecHttpStatus = svh

		cpuMetrics := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "system_cpu_usage",
			Help: "[cpu] cpu usage (percent)",
		}, []string{"core"})
		metricsSummaryGaugeVecCPUUsage = cpuMetrics

		memMetrics := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "system_memory_usage",
			Help: "[memory] memory usage bytes",
		}, []string{"type"})
		metricsSummaryGaugeVecMemoryUsage = memMetrics

		diskMetrics := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "system_disk_usage",
			Help: "[disk] disk usage bytes",
		}, []string{"device", "type"})
		metricsSummaryGaugeVecDiskUsage = diskMetrics

		networkMetrics := prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "system_network_transfer",
			Help: "[network] network transfer bytes",
		}, []string{"name", "type"})
		metricsSummaryGaugeVecNetworkTransfer = networkMetrics

		prometheus.MustRegister(cpuMetrics, memMetrics, diskMetrics, networkMetrics)
	})
}

// SystemMetricsSample samples system metrics
//
// The default sampling interval is once per minute, overwriting the previous values
//
// Collected metrics:
//
// - CPU, multiple cores -> utilization
//
// - Memory -> usage and utilization
//
// - Disk usage, root path -> usage and utilization
//
// - Disk I/O, all devices -> read/write counts, bytes, and time
//
// - Network interfaces, multiple interfaces -> inbound and outbound values
func SystemMetricsSample(exit chan struct{}) {
	ticker := time.NewTicker(sampleInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			//CPU utilization
			cpuUsage, err := cpu.Percent(sampleInterval, true)
			if err == nil && len(cpuUsage) > 0 {
				for index, cu := range cpuUsage {
					metricsSummaryGaugeVecCPUUsage.WithLabelValues(fmt.Sprintf("%d", index)).Set(cu)
				}
			}
			//Memory utilization
			memStat, err := mem.VirtualMemory()
			if err == nil && memStat != nil {
				metricsSummaryGaugeVecMemoryUsage.WithLabelValues("used").Set(float64(memStat.Used))
				metricsSummaryGaugeVecMemoryUsage.WithLabelValues("percent").Set(memStat.UsedPercent)
			}
			//Disk utilization
			diskStat, err := disk.Usage(diskPath)
			if err == nil && diskStat != nil {
				metricsSummaryGaugeVecDiskUsage.WithLabelValues(diskPath, "used").Set(float64(diskStat.Used))
				metricsSummaryGaugeVecDiskUsage.WithLabelValues(diskPath, "percent").Set(diskStat.UsedPercent)
			}
			//Disk I/O
			diskIoStat, err := disk.IOCounters()
			if err == nil && diskIoStat != nil {
				for device, stat := range diskIoStat {
					//Read
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "readBytes").Set(float64(stat.ReadBytes))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "readTime").Set(float64(stat.ReadTime))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "readCount").Set(float64(stat.ReadCount))

					//Write
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "writeBytes").Set(float64(stat.WriteBytes))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "writeTime").Set(float64(stat.WriteTime))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "writeCount").Set(float64(stat.WriteCount))
				}
			}
			//Network interface utilization
			netStats, err := net.IOCounters(false)
			if err == nil && len(netStats) > 0 {
				for _, netStat := range netStats {
					metricsSummaryGaugeVecNetworkTransfer.WithLabelValues(netStat.Name, "bytesReceived").Set(float64(netStat.BytesRecv))
					metricsSummaryGaugeVecNetworkTransfer.WithLabelValues(netStat.Name, "bytesSent").Set(float64(netStat.BytesSent))
				}
			}
		case <-exit:
			break LOOP
		}
	}
}
