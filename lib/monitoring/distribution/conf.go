package distribution

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

const (
	// DefaultLocalWindow is the default window duration for a single node's Sketch.
	DefaultLocalWindow = 10 * time.Second
	// DefaultAggregateWindow is the default window duration for distributed aggregation results.
	DefaultAggregateWindow = time.Minute
	// DefaultAllowedLateness is the default wait for late data before an aggregation window closes.
	DefaultAllowedLateness = 20 * time.Second
	// DefaultRelativeAccuracy is the default upper bound on relative error for DDSketch.
	DefaultRelativeAccuracy = 0.01
	// DefaultRetention is the default retention period for node-level Sketch records.
	DefaultRetention = 24 * time.Hour
	// DefaultReporterQueueSize is the default number of records the reporting queue can buffer.
	DefaultReporterQueueSize = 1024
	// DefaultLeaderTTL is the default lifetime of the distributed aggregation lease.
	DefaultLeaderTTL = 15 * time.Second
)

// Conf controls local recording, reporting, and distributed aggregation behavior.
type Conf struct {
	Enabled   bool   `yaml:"enabled" toml:"enabled" json:"enabled"`       // Enabled controls whether distributed quantile statistics are enabled.
	Namespace string `yaml:"namespace" toml:"namespace" json:"namespace"` // Namespace isolates environments or tenants.
	Service   string `yaml:"service" toml:"service" json:"service"`       // Service identifies the application service producing metrics.
	NodeID    string `yaml:"nodeID" toml:"nodeID" json:"nodeID"`          // NodeID uniquely identifies a service instance.

	LocalWindow      time.Duration `yaml:"localWindow" toml:"localWindow" json:"localWindow"`                // LocalWindow is the window duration for each node-level Sketch.
	AggregateWindow  time.Duration `yaml:"aggregateWindow" toml:"aggregateWindow" json:"aggregateWindow"`    // AggregateWindow is the window duration for merging across nodes.
	AllowedLateness  time.Duration `yaml:"allowedLateness" toml:"allowedLateness" json:"allowedLateness"`    // AllowedLateness delays window closure to accommodate late reports.
	RelativeAccuracy float64       `yaml:"relativeAccuracy" toml:"relativeAccuracy" json:"relativeAccuracy"` // RelativeAccuracy must be between zero and one.
	Quantiles        []float64     `yaml:"quantiles" toml:"quantiles" json:"quantiles"`                      // Quantiles lists the quantiles to calculate; each must be in the open interval (0, 1).

	Reporter   ReporterConfig   `yaml:"reporter" toml:"reporter" json:"reporter"`       // Reporter configures record reporting and retention policies.
	Aggregator AggregatorConfig `yaml:"aggregator" toml:"aggregator" json:"aggregator"` // Aggregator configures leader election and cross-node aggregation behavior.
}

// ReporterConfig controls the reporting backend, queue, and retention policy for node-level Sketches.
type ReporterConfig struct {
	Driver    string        `yaml:"driver" toml:"driver" json:"driver"`          // Driver specifies the reporting backend, such as redis.
	QueueSize int           `yaml:"queueSize" toml:"queueSize" json:"queueSize"` // QueueSize limits the capacity of the nonblocking reporting queue.
	Retention time.Duration `yaml:"retention" toml:"retention" json:"retention"` // Retention controls how long node-level records are retained.
}

// AggregatorConfig controls leader election and aggregated result retention policies.
type AggregatorConfig struct {
	Enabled         bool          `yaml:"enabled" toml:"enabled" json:"enabled"`                         // Enabled controls whether to start the aggregator on the current node.
	LeaderKey       string        `yaml:"leaderKey" toml:"leaderKey" json:"leaderKey"`                   // LeaderKey is the shared lease key for leader election.
	LeaderTTL       time.Duration `yaml:"leaderTTL" toml:"leaderTTL" json:"leaderTTL"`                   // LeaderTTL is the lifetime of an acquired lease.
	ScanInterval    time.Duration `yaml:"scanInterval" toml:"scanInterval" json:"scanInterval"`          // ScanInterval controls the interval for lease renewal and window scanning.
	ResultRetention time.Duration `yaml:"resultRetention" toml:"resultRetention" json:"resultRetention"` // ResultRetention controls how long aggregated results are retained.
}

// SetDefaults fills in unset configuration defaults and sorts Quantiles.
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

// Validate checks that the configuration is complete and internally consistent when enabled.
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
