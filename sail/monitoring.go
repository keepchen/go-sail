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

// Monitor 提供分布式统计初始化和聚合结果查询能力。
type Monitor interface {
	// Distribution 初始化分布式统计组件。同一个监控实例只会执行一次有效初始化。
	// disabled 配置不会占用初始化机会，后续仍可使用 enabled 配置完成初始化。
	Distribution(conf *distribution.Conf, redisClient redis.UniversalClient)
	// Result 获取指定指标在聚合窗口
	//
	// [windowStart, windowEnd) 内的完整统计结果。
	Result(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (*distribution.AggregatedResult, error)
	// P 按 Pxx 记法获取指定指标的分位值。
	//
	// bool 表示聚合结果中是否包含该分位点。
	P(percentile uint, metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P85 获取指定指标的 P85 分位值。
	P85(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P99 获取指定指标的 P99 分位值。
	P99(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P990 获取指定指标的 P99.0 分位值，与 P99 数值等价。
	P990(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error)
	// P999 获取指定指标的 P99.9 分位值。
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

// Monitoring 初始化并返回全局唯一的监控实例。
// 首次调用传入的 ctx 和 logger 会用于该实例的整个生命周期，后续调用仅返回已有实例。
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

// GetMonitor 返回已由 Monitoring 初始化的全局监控实例。
// 尚未初始化时返回 nil。
func GetMonitor() Monitor {
	monitorMu.RLock()
	defer monitorMu.RUnlock()

	return monitorInstance
}

// Distribution 初始化分布式统计组件。同一个监控实例只会执行一次有效初始化。
// disabled 配置不会占用初始化机会，后续仍可使用 enabled 配置完成初始化。
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
		// 如果当前节点开启了聚合，就执行聚合操作。
		if conf.Aggregator.Enabled {
			aggregator := distribution.NewAggregator(conf, redisStore, leadership, m.logger)
			go aggregator.Run(m.ctx)
		}

		// 连接到 HTTP 埋点。
		api.SetRegistry(registry)
	})
}

// Result 获取指定指标在聚合窗口
//
// [windowStart, windowEnd) 内的完整统计结果。
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

// P 按 Pxx 记法获取指定指标的分位值。
//
// bool 表示聚合结果中是否包含该分位点。
func (m *monitoringImpl) P(percentile uint, metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	result, err := m.Result(metric, labels, windowStart, windowEnd)
	if err != nil {
		return 0, false, err
	}
	value, ok := result.P(percentile)
	return value, ok, nil
}

// P85 获取指定指标的 P85 分位值。
func (m *monitoringImpl) P85(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(85, metric, labels, windowStart, windowEnd)
}

// P99 获取指定指标的 P99 分位值。
func (m *monitoringImpl) P99(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(99, metric, labels, windowStart, windowEnd)
}

// P990 获取指定指标的 P99.0 分位值，与 P99 数值等价。
func (m *monitoringImpl) P990(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(990, metric, labels, windowStart, windowEnd)
}

// P999 获取指定指标的 P99.9 分位值。
func (m *monitoringImpl) P999(metric string, labels distribution.Labels, windowStart, windowEnd time.Time) (float64, bool, error) {
	return m.P(999, metric, labels, windowStart, windowEnd)
}
