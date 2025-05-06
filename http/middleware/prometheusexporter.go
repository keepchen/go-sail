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

// PrometheusExporter 导出Prometheus指标
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
	//记录接口延时
	metricsSummaryVecHttpLatency.WithLabelValues(r.ctx.FullPath()).Observe(float64(elapsed))
	//记录接口返回状态码
	metricsSummaryVecHttpStatus.WithLabelValues(r.ctx.FullPath(), fmt.Sprintf("%d", r.ctx.Writer.Status())).Observe(float64(1))

	r.written = true

	return r.ResponseWriter.Write(b)
}

var once sync.Once

// http接口状态与延迟
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

// SetDiskPath 设置磁盘监控路径
func SetDiskPath(path string) {
	if len(path) != 0 {
		diskPath = path
	}
}

// SetSampleInterval 设置采样间隔(频率)
func SetSampleInterval(interval string) {
	td, err := time.ParseDuration(interval)
	//若小于1ms，则使用默认值1分钟
	if err != nil || td.Milliseconds() == 0 {
		return
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

// SystemMetricsSample 系统指标采样
//
// 默认采样周期为每分钟一次，覆盖上一次记录值
//
// 采集对象：
//
// - cpu，多核 -> 使用率
//
// - 内存 -> 使用数值和使用率
//
// - 硬盘使用率，根路径 -> 使用数值和使用率
//
// - 硬盘io，遍历设备 -> 读写次数、字节数、时间
//
// - 网卡，多网卡 -> 出入站数值
func SystemMetricsSample(exit chan struct{}) {
	ticker := time.NewTicker(sampleInterval)

LOOP:
	for {
		select {
		case <-ticker.C:
			//cpu使用率
			cpuUsage, err := cpu.Percent(sampleInterval, true)
			if err == nil && len(cpuUsage) > 0 {
				for index, cu := range cpuUsage {
					metricsSummaryGaugeVecCPUUsage.WithLabelValues(fmt.Sprintf("%d", index)).Set(cu)
				}
			}
			//内存使用率
			memStat, err := mem.VirtualMemory()
			if err == nil && memStat != nil {
				metricsSummaryGaugeVecMemoryUsage.WithLabelValues("used").Set(float64(memStat.Used))
				metricsSummaryGaugeVecMemoryUsage.WithLabelValues("percent").Set(memStat.UsedPercent)
			}
			//硬盘使用率
			diskStat, err := disk.Usage(diskPath)
			if err == nil && diskStat != nil {
				metricsSummaryGaugeVecDiskUsage.WithLabelValues(diskPath, "used").Set(float64(diskStat.Used))
				metricsSummaryGaugeVecDiskUsage.WithLabelValues(diskPath, "percent").Set(diskStat.UsedPercent)
			}
			//硬盘io
			diskIoStat, err := disk.IOCounters()
			if err == nil && diskIoStat != nil {
				for device, stat := range diskIoStat {
					//读
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "readBytes").Set(float64(stat.ReadBytes))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "readTime").Set(float64(stat.ReadTime))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "readCount").Set(float64(stat.ReadCount))

					//写
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "writeBytes").Set(float64(stat.WriteBytes))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "writeTime").Set(float64(stat.WriteTime))
					metricsSummaryGaugeVecDiskUsage.WithLabelValues(device, "writeCount").Set(float64(stat.WriteCount))
				}
			}
			//网卡使用率
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
