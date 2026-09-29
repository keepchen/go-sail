package distribution

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

const (
	// DefaultLocalWindow 表示单节点生成 Sketch 的默认窗口时长。
	DefaultLocalWindow = 10 * time.Second
	// DefaultAggregateWindow 表示分布式聚合结果的默认窗口时长。
	DefaultAggregateWindow = time.Minute
	// DefaultAllowedLateness 表示聚合窗口关闭前等待迟到数据的默认时长。
	DefaultAllowedLateness = 20 * time.Second
	// DefaultRelativeAccuracy 表示 DDSketch 默认的相对误差上限。
	DefaultRelativeAccuracy = 0.01
	// DefaultRetention 表示节点级 Sketch 记录的默认保留时长。
	DefaultRetention = 24 * time.Hour
	// DefaultReporterQueueSize 表示上报队列默认可缓冲的记录数。
	DefaultReporterQueueSize = 1024
	// DefaultLeaderTTL 表示分布式聚合租约的默认有效期。
	DefaultLeaderTTL = 15 * time.Second
)

// Conf 控制本地记录、数据上报和分布式聚合行为。
type Conf struct {
	Enabled   bool   `yaml:"enabled" toml:"enabled" json:"enabled"`       // Enabled 控制是否启用分布式分位数统计。
	Namespace string `yaml:"namespace" toml:"namespace" json:"namespace"` // Namespace 用于隔离环境或租户。
	Service   string `yaml:"service" toml:"service" json:"service"`       // Service 标识产生指标的应用服务。
	NodeID    string `yaml:"nodeID" toml:"nodeID" json:"nodeID"`          // NodeID 唯一标识一个服务实例。

	LocalWindow      time.Duration `yaml:"localWindow" toml:"localWindow" json:"localWindow"`                // LocalWindow 表示每个节点级 Sketch 的窗口时长。
	AggregateWindow  time.Duration `yaml:"aggregateWindow" toml:"aggregateWindow" json:"aggregateWindow"`    // AggregateWindow 表示跨节点合并的窗口时长。
	AllowedLateness  time.Duration `yaml:"allowedLateness" toml:"allowedLateness" json:"allowedLateness"`    // AllowedLateness 为迟到上报延后窗口关闭时间。
	RelativeAccuracy float64       `yaml:"relativeAccuracy" toml:"relativeAccuracy" json:"relativeAccuracy"` // RelativeAccuracy 必须位于零和一之间。
	Quantiles        []float64     `yaml:"quantiles" toml:"quantiles" json:"quantiles"`                      // Quantiles 是待计算的分位点，每项必须位于开区间 (0, 1)。

	Reporter   ReporterConfig   `yaml:"reporter" toml:"reporter" json:"reporter"`       // Reporter 配置记录上报和保留策略。
	Aggregator AggregatorConfig `yaml:"aggregator" toml:"aggregator" json:"aggregator"` // Aggregator 配置选主和跨节点聚合行为。
}

// ReporterConfig 控制节点级 Sketch 的上报后端、队列和保留策略。
type ReporterConfig struct {
	Driver    string        `yaml:"driver" toml:"driver" json:"driver"`          // Driver 指定上报后端，例如 redis。
	QueueSize int           `yaml:"queueSize" toml:"queueSize" json:"queueSize"` // QueueSize 限制非阻塞上报队列的容量。
	Retention time.Duration `yaml:"retention" toml:"retention" json:"retention"` // Retention 控制节点级记录的保留时长。
}

// AggregatorConfig 控制选主和聚合结果的保留策略。
type AggregatorConfig struct {
	Enabled         bool          `yaml:"enabled" toml:"enabled" json:"enabled"`                         // Enabled 控制是否在当前节点启动聚合器。
	LeaderKey       string        `yaml:"leaderKey" toml:"leaderKey" json:"leaderKey"`                   // LeaderKey 是参与选主的共享租约键。
	LeaderTTL       time.Duration `yaml:"leaderTTL" toml:"leaderTTL" json:"leaderTTL"`                   // LeaderTTL 是已获取租约的有效期。
	ScanInterval    time.Duration `yaml:"scanInterval" toml:"scanInterval" json:"scanInterval"`          // ScanInterval 控制租约续期和窗口扫描间隔。
	ResultRetention time.Duration `yaml:"resultRetention" toml:"resultRetention" json:"resultRetention"` // ResultRetention 控制聚合结果的保留时长。
}

// SetDefaults 填充未设置的配置默认值，并对 Quantiles 进行排序。
func (c *Conf) SetDefaults() {
	if c.LocalWindow <= 0 {
		c.LocalWindow = DefaultLocalWindow
	}
	if c.AggregateWindow <= 0 {
		c.AggregateWindow = DefaultAggregateWindow
	}
	if c.AllowedLateness <= 0 {
		c.AllowedLateness = DefaultAllowedLateness
	}
	if c.RelativeAccuracy <= 0 {
		c.RelativeAccuracy = DefaultRelativeAccuracy
	}
	if len(c.Quantiles) == 0 {
		c.Quantiles = []float64{0.85, 0.99, 0.999}
	}
	if c.Reporter.QueueSize <= 0 {
		c.Reporter.QueueSize = DefaultReporterQueueSize
	}
	if c.Reporter.Retention <= 0 {
		c.Reporter.Retention = DefaultRetention
	}
	if c.Aggregator.LeaderTTL <= 0 {
		c.Aggregator.LeaderTTL = DefaultLeaderTTL
	}
	if c.Aggregator.ScanInterval <= 0 {
		c.Aggregator.ScanInterval = 5 * time.Second
	}
	if c.Aggregator.ResultRetention <= 0 {
		c.Aggregator.ResultRetention = 7 * 24 * time.Hour
	}
	if c.Aggregator.LeaderKey == "" {
		c.Aggregator.LeaderKey = "go-sail:distribution:aggregator:leader"
	}

	sort.Float64s(c.Quantiles)
}

// Validate 检查启用状态下的配置是否完整且内部一致。
func (c *Conf) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Namespace == "" {
		return errors.New("[Go-Sail] <monitoring> distribution namespace is required")
	}
	if c.Service == "" {
		return errors.New("[Go-Sail] <monitoring> distribution service is required")
	}
	if c.NodeID == "" {
		return errors.New("[Go-Sail] <monitoring> distribution node ID is required")
	}
	if c.LocalWindow <= 0 {
		return errors.New("[Go-Sail] <monitoring> local window must be greater than zero")
	}
	if c.AggregateWindow <= 0 {
		return errors.New("[Go-Sail] <monitoring> aggregate window must be greater than zero")
	}
	if c.AggregateWindow%c.LocalWindow != 0 {
		return errors.New("[Go-Sail] <monitoring> aggregate window must be divisible by local window")
	}
	if c.RelativeAccuracy <= 0 || c.RelativeAccuracy >= 1 {
		return errors.New("[Go-Sail] <monitoring> relative accuracy must be between zero and one")
	}
	for _, q := range c.Quantiles {
		if q <= 0 || q >= 1 {
			return fmt.Errorf("[Go-Sail] <monitoring> invalid quantile: %v", q)
		}
	}
	if c.Aggregator.LeaderTTL <= c.Aggregator.ScanInterval {
		return errors.New("[Go-Sail] <monitoring> leader TTL must be greater than scan interval")
	}
	return nil
}
