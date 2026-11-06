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
	// SchemaVersionV1 identifies the first version of the record and result transport schema.
	SchemaVersionV1 uint16 = 1
	// AlgorithmDDSketchV1 identifies the DDSketch payload format used by the first schema version.
	AlgorithmDDSketchV1 = "ddsketch-v1"
	// UnitMicrosecond indicates that values are measured in microseconds.
	UnitMicrosecond = "microsecond"
)

// SketchRecord represents a mergeable Sketch and exact summary produced by one node within a local window.
// WindowStart is inclusive and WindowEnd is exclusive.
type SketchRecord struct {
	SchemaVersion uint16 `json:"schemaVersion"` // Transport schema version.
	Algorithm     string `json:"algorithm"`     // Algorithm and version of the Sketch payload.

	RecordID  string `json:"recordId"`  // Stable record identifier for idempotent reporting.
	MetricKey string `json:"metricKey"` // Stable identifier for the metric series.

	Namespace string `json:"namespace"` // Namespace or tenant the metric belongs to.
	Service   string `json:"service"`   // Service producing the metric.
	Metric    string `json:"metric"`    // Metric name.
	Labels    Labels `json:"labels"`    // Metric dimensions.

	NodeID string `json:"nodeId"` // Service instance producing the record.

	WindowStart time.Time `json:"windowStart"` // Inclusive boundary of the local window.
	WindowEnd   time.Time `json:"windowEnd"`   // Exclusive boundary of the local window.

	Unit             string  `json:"unit"`             // Unit shared by exact statistics and Sketch samples.
	RelativeAccuracy float64 `json:"relativeAccuracy"` // Upper bound on relative error used to build the Sketch.

	Count   uint64  `json:"count"`   // Exact sample count.
	Sum     float64 `json:"sum"`     // Exact sample sum.
	Min     float64 `json:"min"`     // Exact minimum sample value.
	Max     float64 `json:"max"`     // Exact maximum sample value.
	Payload []byte  `json:"payload"` // Encoded mergeable Sketch including mapping information.
}

// BuildRecordID generates a stable record identifier from the schema version, metric series, producer node, and window boundaries.
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

// AggregatedResult represents distributed statistics for a metric series within an aggregation window.
// WindowStart is inclusive and WindowEnd is exclusive.
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

	NodeCount   uint64 `json:"nodeCount"`   // Number of distinct nodes included in aggregation.
	RecordCount uint64 `json:"recordCount"` // Number of distinct local records merged.
	SampleCount uint64 `json:"sampleCount"` // Exact total sample count.

	Sum float64 `json:"sum"`
	Min float64 `json:"min"`
	Max float64 `json:"max"`
	Avg float64 `json:"avg"`

	Quantiles map[string]float64 `json:"quantiles"` // Approximate quantile values keyed by formatted quantiles.

	Payload   []byte    `json:"payload"`   // Aggregated Sketch for further downstream merging.
	Finalized bool      `json:"finalized"` // Whether the allowed lateness has elapsed and the result is closed.
	CreatedAt time.Time `json:"createdAt"` // UTC creation time.
}

// Quantile returns the aggregated value for a quantile. quantile must be in the open interval (0, 1).
// The returned bool indicates whether the quantile is included in the aggregated result.
func (r *AggregatedResult) Quantile(quantile float64) (float64, bool) {
	if r == nil || quantile <= 0 || quantile >= 1 || math.IsNaN(quantile) {
		return 0, false
	}

	value, ok := r.Quantiles[formatQuantile(quantile)]
	return value, ok
}

// P returns a quantile value in Pxx notation. Each digit of percentile is treated as a fractional digit,
// e.g. P(85), P(99), and P(999) represent P85, P99, and P999, respectively.
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

// P85 returns the P85 quantile value.
func (r *AggregatedResult) P85() (float64, bool) {
	return r.P(85)
}

// P99 returns the P99 quantile value.
func (r *AggregatedResult) P99() (float64, bool) {
	return r.P(99)
}

// P990 returns the P99.0 quantile value, numerically equivalent to P99.
func (r *AggregatedResult) P990() (float64, bool) {
	return r.P(990)
}

// P999 returns the P99.9 quantile value.
func (r *AggregatedResult) P999() (float64, bool) {
	return r.P(999)
}

// BuildResultID returns a stable result identifier for a metric aggregation window.
func BuildResultID(metricKey string, windowStart time.Time, windowEnd time.Time) string {
	raw := strings.Join([]string{
		metricKey,
		strconv.FormatInt(windowStart.UnixNano(), 10),
		strconv.FormatInt(windowEnd.UnixNano(), 10),
	}, "\x00")

	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
