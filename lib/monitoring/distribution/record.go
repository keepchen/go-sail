package distribution

import (
	"crypto/sha256"
	"encoding/hex"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	// SchemaVersionV1 标识记录和结果的第一版传输结构。
	SchemaVersionV1 uint16 = 1
	// AlgorithmDDSketchV1 标识第一版结构使用的 DDSketch 载荷格式。
	AlgorithmDDSketchV1 = "ddsketch-v1"
	// UnitMicrosecond 表示数值采用微秒作为计量单位。
	UnitMicrosecond = "microsecond"
)

// SketchRecord 表示单个节点在一个本地窗口内生成的可合并 Sketch 和精确摘要。
// WindowStart 为闭区间边界，WindowEnd 为开区间边界。
type SketchRecord struct {
	SchemaVersion uint16 `json:"schemaVersion"` // 传输结构版本。
	Algorithm     string `json:"algorithm"`     // Sketch 载荷的算法及版本。

	RecordID  string `json:"recordId"`  // 用于幂等上报的稳定记录标识。
	MetricKey string `json:"metricKey"` // 指标序列的稳定标识。

	Namespace string `json:"namespace"` // 指标所属命名空间或租户。
	Service   string `json:"service"`   // 产生指标的服务。
	Metric    string `json:"metric"`    // 指标名称。
	Labels    Labels `json:"labels"`    // 指标维度。

	NodeID string `json:"nodeId"` // 产生记录的服务实例。

	WindowStart time.Time `json:"windowStart"` // 本地窗口的闭区间边界。
	WindowEnd   time.Time `json:"windowEnd"`   // 本地窗口的开区间边界。

	Unit             string  `json:"unit"`             // 精确统计值和 Sketch 样本共用的单位。
	RelativeAccuracy float64 `json:"relativeAccuracy"` // 构建 Sketch 时采用的相对误差上限。

	Count   uint64  `json:"count"`   // 精确样本数。
	Sum     float64 `json:"sum"`     // 精确样本总和。
	Min     float64 `json:"min"`     // 精确最小样本值。
	Max     float64 `json:"max"`     // 精确最大样本值。
	Payload []byte  `json:"payload"` // 包含映射信息的可合并 Sketch 编码。
}

// BuildRecordID 根据结构版本、指标序列、生产节点和窗口边界生成稳定记录标识。
func (r *SketchRecord) BuildRecordID() string {
	raw := strings.Join([]string{
		strconv.Itoa(int(r.SchemaVersion)),
		r.Algorithm,
		r.MetricKey,
		r.NodeID,
		strconv.FormatInt(r.WindowStart.UnixNano(), 10),
		strconv.FormatInt(r.WindowEnd.UnixNano(), 10),
	}, "\x00")

	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// AggregatedResult 表示一条指标序列在一个聚合窗口内的分布式统计结果。
// WindowStart 为闭区间边界，WindowEnd 为开区间边界。
type AggregatedResult struct {
	SchemaVersion uint16 `json:"schemaVersion"`
	Algorithm     string `json:"algorithm"`

	ResultID  string `json:"resultId"`
	MetricKey string `json:"metricKey"`

	Namespace string `json:"namespace"`
	Service   string `json:"service"`
	Metric    string `json:"metric"`
	Labels    Labels `json:"labels"`

	WindowStart time.Time `json:"windowStart"`
	WindowEnd   time.Time `json:"windowEnd"`

	RelativeAccuracy float64 `json:"relativeAccuracy"`

	NodeCount   uint64 `json:"nodeCount"`   // 参与聚合的不同节点数。
	RecordCount uint64 `json:"recordCount"` // 合并的不同本地记录数。
	SampleCount uint64 `json:"sampleCount"` // 精确样本总数。

	Sum float64 `json:"sum"`
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`

	Quantiles map[string]float64 `json:"quantiles"` // 以格式化分位点为键的近似分位值。

	Payload   []byte    `json:"payload"`   // 供下游继续合并的聚合 Sketch。
	Finalized bool      `json:"finalized"` // 是否已超过迟到等待期并关闭结果。
	CreatedAt time.Time `json:"createdAt"` // UTC 创建时间。
}

// Quantile 返回指定分位点的聚合值。quantile 必须位于开区间 (0, 1)。
// 返回值中的 bool 表示该分位点是否包含在聚合结果中。
func (r *AggregatedResult) Quantile(quantile float64) (float64, bool) {
	if r == nil || quantile <= 0 || quantile >= 1 || math.IsNaN(quantile) {
		return 0, false
	}

	value, ok := r.Quantiles[formatQuantile(quantile)]
	return value, ok
}

// P 按 Pxx 记法返回分位值。percentile 的每一位均视为小数位，
// 例如 P(85)、P(99)、P(999) 分别表示 P85、P99、P999。
func (r *AggregatedResult) P(percentile uint) (float64, bool) {
	if percentile == 0 {
		return 0, false
	}

	digits := 1
	for value := percentile; value >= 10; value /= 10 {
		digits++
	}
	quantile := float64(percentile) / math.Pow10(digits)

	return r.Quantile(quantile)
}

// P85 返回 P85 分位值。
func (r *AggregatedResult) P85() (float64, bool) {
	return r.P(85)
}

// P99 返回 P99 分位值。
func (r *AggregatedResult) P99() (float64, bool) {
	return r.P(99)
}

// P990 返回 P99.0 分位值，与 P99 数值等价。
func (r *AggregatedResult) P990() (float64, bool) {
	return r.P(990)
}

// P999 返回 P99.9 分位值。
func (r *AggregatedResult) P999() (float64, bool) {
	return r.P(999)
}

// BuildResultID 返回指标聚合窗口的稳定结果标识。
func BuildResultID(metricKey string, windowStart time.Time, windowEnd time.Time) string {
	raw := strings.Join([]string{
		metricKey,
		strconv.FormatInt(windowStart.UnixNano(), 10),
		strconv.FormatInt(windowEnd.UnixNano(), 10),
	}, "\x00")

	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
