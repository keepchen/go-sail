package sail

import (
	"context"
	"testing"
	"time"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
	"go.uber.org/zap"
)

type monitoringResultStore struct {
	resultID string
	result   *distribution.AggregatedResult
}

func (s *monitoringResultStore) Report(context.Context, *distribution.SketchRecord) error {
	return nil
}
func (s *monitoringResultStore) ListRecords(context.Context, distribution.WindowQuery) ([]*distribution.SketchRecord, error) {
	return nil, nil
}
func (s *monitoringResultStore) SaveResult(context.Context, *distribution.AggregatedResult) error {
	return nil
}
func (s *monitoringResultStore) GetResult(_ context.Context, resultID string) (*distribution.AggregatedResult, error) {
	s.resultID = resultID
	return s.result, nil
}
func (s *monitoringResultStore) MarkAggregated(context.Context, []string) error { return nil }
func (s *monitoringResultStore) GetAggregationCursor(context.Context, string, string) (time.Time, error) {
	return time.Time{}, nil
}
func (s *monitoringResultStore) SaveAggregationCursor(context.Context, string, string, time.Time) error {
	return nil
}

func TestMonitoringResultSugar(t *testing.T) {
	windowStart := time.Date(2026, 8, 12, 10, 0, 0, 0, time.FixedZone("CST", 8*60*60))
	windowEnd := windowStart.Add(time.Minute)
	labels := distribution.Labels{"method": "GET", "route": "/users", "status": "2xx"}
	store := &monitoringResultStore{result: &distribution.AggregatedResult{
		Quantiles: map[string]float64{"0.85": 850, "0.99": 990, "0.999": 999},
	}}
	monitor := &monitoringImpl{
		ctx:    context.Background(),
		logger: zap.NewNop(),
		conf:   &distribution.Conf{Namespace: "production", Service: "api"},
		store:  store,
	}

	value, ok, err := monitor.P990("http.server.duration", labels, windowStart, windowEnd)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || value != 990 {
		t.Fatalf("P990 = (%v, %v)，期望 (990, true)", value, ok)
	}

	metricKey := distribution.BuildMetricKey("production", "api", "http.server.duration", labels)
	wantResultID := distribution.BuildResultID(metricKey, windowStart.UTC(), windowEnd.UTC())
	if store.resultID != wantResultID {
		t.Fatalf("ResultID = %s，期望 %s", store.resultID, wantResultID)
	}
}

func TestMonitoringResultBeforeInitialization(t *testing.T) {
	monitor := &monitoringImpl{ctx: context.Background(), logger: zap.NewNop()}
	if _, err := monitor.Result("metric", nil, time.Now(), time.Now().Add(time.Minute)); err == nil {
		t.Fatal("初始化前查询应返回错误")
	}
}

func TestMonitoringSingleton(t *testing.T) {
	first := Monitoring(context.Background(), zap.NewNop())
	second := Monitoring(context.TODO(), zap.NewExample())

	if first == nil {
		t.Fatal("Monitoring 首次初始化不应返回 nil")
	}
	if first != second {
		t.Fatal("Monitoring 重复调用应返回同一个实例")
	}
	if GetMonitor() != first {
		t.Fatal("GetMonitor 应返回已初始化的单例")
	}
}

func TestMonitoringDistributionOnlyInitializesOnce(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	monitor := &monitoringImpl{ctx: ctx, logger: zap.NewNop()}

	disabled := &distribution.Conf{Enabled: false}
	monitor.Distribution(disabled, nil)
	if monitor.conf != nil || monitor.store != nil {
		t.Fatal("disabled 配置不应初始化 Distribution")
	}

	first := &distribution.Conf{
		Enabled:   true,
		Namespace: "production",
		Service:   "api",
		NodeID:    "node-1",
	}
	monitor.Distribution(first, nil)

	second := &distribution.Conf{
		Enabled:   true,
		Namespace: "staging",
		Service:   "other-api",
		NodeID:    "node-2",
	}
	monitor.Distribution(second, nil)
	cancel()

	if monitor.conf != first {
		t.Fatal("Distribution 重复初始化不应覆盖首次配置")
	}
}
