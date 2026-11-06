package sail

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/http/api"
	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
	redisbackend "github.com/keepchen/go-sail/v3/lib/monitoring/distribution/backend/redis"
	"go.uber.org/zap"
)

// Monitor provides distributed statistics initialization and aggregated result queries.
type Monitor interface {
	// Distribution initializes distributed statistics. Each monitor instance performs successful initialization only once.
	// A disabled configuration does not consume the initialization opportunity; a later enabled configuration can still initialize it.
	Distribution(conf *distribution.Conf, redisClient redis.UniversalClient)
	// Result retrieves the complete statistics for the specified metric within the aggregation window
	//
	// [windowStart, windowEnd).
	Result(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (*distribution.AggregatedResult, error)
	// P retrieves a quantile value for the specified metric using Pxx notation.
	//
	// bool indicates whether the aggregated result contains the quantile.
	P(percentile uint, metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P85 retrieves the P85 quantile value for the specified metric.
	P85(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P99 retrieves the P99 quantile value for the specified metric.
	P99(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P990 retrieves the P99.0 quantile value for the specified metric, numerically equivalent to P99.
	P990(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P999 retrieves the P99.9 quantile value for the specified metric.
	P999(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
}

type monitoringImpl struct {
	ctx    context.Context
	logger *zap.Logger

	initOnce sync.Once
	mu       sync.RWMutex
	conf     *distribution.Conf
	store    distribution.Store
}

var _ Monitor = (*monitoringImpl)(nil)

var (
	monitorOnce     sync.Once
	monitorMu       sync.RWMutex
	monitorInstance Monitor
)

// Monitoring initializes and returns the single global monitor instance.
// The ctx and logger from the first call are used throughout the instance's lifetime; subsequent calls only return the existing instance.
func Monitoring(ctx context.Context, logger *zap.Logger) Monitor {
	monitorOnce.Do(func() {
		monitorMu.Lock()
		monitorInstance = &monitoringImpl{
			ctx:    ctx,
			logger: logger,
		}
		monitorMu.Unlock()
	})

	return GetMonitor()
}

// GetMonitor returns the global monitor instance initialized by Monitoring.
// Returns nil if not yet initialized.
func GetMonitor() Monitor {
	monitorMu.RLock()
	defer monitorMu.RUnlock()

	return monitorInstance
}

// Distribution initializes distributed statistics. Each monitor instance performs successful initialization only once.
// A disabled configuration does not consume the initialization opportunity; a later enabled configuration can still initialize it.
func (m *monitoringImpl) Distribution(conf *distribution.Conf, redisClient redis.UniversalClient) {
	if !conf.Enabled {
		return
	}

	m.initOnce.Do(func() {
		registry, err := distribution.NewRegistry(conf)
		if err != nil {
			panic(err)
		}
		redisStore := redisbackend.NewStore(redisClient, conf.LocalWindow, conf.Reporter.Retention, conf.Aggregator.ResultRetention)
		leadership, err := redisbackend.NewLeadership(redisClient, conf.Aggregator.LeaderKey, conf.Aggregator.LeaderTTL)
		if err != nil {
			panic(err)
		}
		reportWorker := distribution.NewReportWorker(redisStore, registry.ReportQueue(), m.logger)

		m.mu.Lock()
		m.conf = conf
		m.store = redisStore
		m.mu.Unlock()

		go registry.Start(m.ctx)
		go reportWorker.Run(m.ctx)
		// Run aggregation if enabled on the current node.
		if conf.Aggregator.Enabled {
			aggregator := distribution.NewAggregator(conf, redisStore, leadership, m.logger)
			go aggregator.Run(m.ctx)
		}

		// Connect to HTTP instrumentation.
		api.SetRegistry(registry)
	})
}

// Result retrieves the complete statistics for the specified metric within the aggregation window
//
// [windowStart, windowEnd).
func (m *monitoringImpl) Result(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (*distribution.AggregatedResult, error) {
	m.mu.RLock()
	conf := m.conf
	store := m.store
	m.mu.RUnlock()

	if conf == nil || store == nil {
		return nil, errors.New("[Go-Sail] <monitoring> distribution is not initialized")
	}
	if metric == "" {
		return nil, errors.New("[Go-Sail] <monitoring> metric is required")
	}
	if !windowStart.Before(windowEnd) {
		return nil, errors.New("[Go-Sail] <monitoring> window start must be before window end")
	}

	metricKey := distribution.BuildMetricKey(conf.Namespace, conf.Service, metric, labels)
	resultID := distribution.BuildResultID(metricKey, windowStart.UTC(), windowEnd.UTC())
	result, err := store.GetResult(m.ctx, resultID)
	if err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> get aggregate result: %w", err)
	}
	return result, nil
}

// P retrieves a quantile value for the specified metric using Pxx notation.
//
// bool indicates whether the aggregated result contains the quantile.
func (m *monitoringImpl) P(percentile uint, metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	result, err := m.Result(metric, labels, windowStart, windowEnd)
	if err != nil {
		return 0, false, err
	}
	value, ok := result.P(percentile)
	return value, ok, nil
}

// P85 retrieves the P85 quantile value for the specified metric.
func (m *monitoringImpl) P85(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(85, metric, labels, windowStart, windowEnd)
}

// P99 retrieves the P99 quantile value for the specified metric.
func (m *monitoringImpl) P99(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(99, metric, labels, windowStart, windowEnd)
}

// P990 retrieves the P99.0 quantile value for the specified metric, numerically equivalent to P99.
func (m *monitoringImpl) P990(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(990, metric, labels, windowStart, windowEnd)
}

// P999 retrieves the P99.9 quantile value for the specified metric.
func (m *monitoringImpl) P999(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(999, metric, labels, windowStart, windowEnd)
}
