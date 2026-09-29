package distribution

import "testing"

func TestAggregatedResultPercentileHelpers(t *testing.T) {
	result := &AggregatedResult{
		Quantiles: map[string]float64{
			"0.85":  85,
			"0.99":  99,
			"0.999": 999,
		},
	}

	tests := []struct {
		name string
		get  func() (float64, bool)
		want float64
	}{
		{name: "Quantile", get: func() (float64, bool) { return result.Quantile(0.85) }, want: 85},
		{name: "P85", get: result.P85, want: 85},
		{name: "P99", get: result.P99, want: 99},
		{name: "P999", get: result.P999, want: 999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.get()
			if !ok || got != tt.want {
				t.Fatalf("获取结果 = (%v, %v)，期望 (%v, true)", got, ok, tt.want)
			}
		})
	}
}

func TestAggregatedResultPercentileHelpersMissing(t *testing.T) {
	result := &AggregatedResult{Quantiles: map[string]float64{"0.99": 99}}

	if _, ok := result.P(95); ok {
		t.Fatal("未配置的 P95 不应存在")
	}
	if _, ok := result.Quantile(1); ok {
		t.Fatal("无效分位点不应存在")
	}
	var nilResult *AggregatedResult
	if _, ok := nilResult.P99(); ok {
		t.Fatal("空结果不应返回分位值")
	}
}
