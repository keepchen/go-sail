package distribution

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution/codec"
	"go.uber.org/zap"
)

// Aggregator 通过选主确定唯一进程，将节点级 Sketch 合并为服务级终态聚合窗口。
type Aggregator struct {
	config *Conf

	store      Store
	leadership Leadership
	logger     *zap.Logger
}

// NewAggregator 使用指定的存储和租约后端创建聚合器。
func NewAggregator(config *Conf, store Store, leadership Leadership, logger *zap.Logger) *Aggregator {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Aggregator{
		config:     config,
		store:      store,
		leadership: leadership,
		logger:     logger,
	}
}

// Run 参与选主并持续聚合已关闭窗口，直至 ctx 被取消。
// 运行期间遗漏的窗口会从持久化游标处依次补算。
func (a *Aggregator) Run(ctx context.Context) {
	if !a.config.Aggregator.Enabled {
		return
	}

	ticker := time.NewTicker(a.config.Aggregator.ScanInterval)
	defer ticker.Stop()

	isLeader := false

	for {
		select {
		case <-ctx.Done():
			if isLeader {
				_ = a.leadership.Release(context.Background())
			}
			return

		case now := <-ticker.C:
			if !isLeader {
				acquired, err := a.leadership.Acquire(ctx)
				if err != nil {
					a.logger.Error(
						"[Go-Sail] <monitoring> acquire distribution aggregator leadership failed",
						zap.Error(err),
					)
					continue
				}

				isLeader = acquired
				if !isLeader {
					continue
				}
			} else {
				renewed, err := a.leadership.Renew(ctx)
				if err != nil || !renewed {
					isLeader = false
					continue
				}
			}

			if err := a.aggregateClosedWindow(ctx, now); err != nil {
				a.logger.Error(
					"[Go-Sail] <monitoring> aggregate distribution window failed",
					zap.Error(err),
				)
			}
		}
	}
}

func (a *Aggregator) aggregateClosedWindow(ctx context.Context, now time.Time) error {
	cutoff := now.Add(-a.config.AllowedLateness)
	latestWindowEnd := cutoff.Truncate(a.config.AggregateWindow)

	if latestWindowEnd.IsZero() {
		return nil
	}

	cursor, err := a.store.GetAggregationCursor(ctx, a.config.Namespace, a.config.Service)
	if err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> load aggregation cursor: %w", err)
	}
	if cursor.IsZero() {
		// 首次启动从最新的已关闭窗口开始；之后由持久化游标补算停机期间的窗口。
		cursor = latestWindowEnd.Add(-a.config.AggregateWindow)
	}

	for windowEnd := cursor.Add(a.config.AggregateWindow); !windowEnd.After(latestWindowEnd); windowEnd = windowEnd.Add(a.config.AggregateWindow) {
		if err := a.AggregateWindow(ctx, windowEnd.Add(-a.config.AggregateWindow), windowEnd); err != nil {
			return err
		}
		if err := a.store.SaveAggregationCursor(ctx, a.config.Namespace, a.config.Service, windowEnd); err != nil {
			return fmt.Errorf("[Go-Sail] <monitoring> save aggregation cursor: %w", err)
		}
	}

	return nil
}

// AggregateWindow 校验、分组并合并完全包含在半开区间
// [windowStart, windowEnd) 内的所有记录。
func (a *Aggregator) AggregateWindow(ctx context.Context, windowStart time.Time, windowEnd time.Time) error {
	records, err := a.store.ListRecords(
		ctx,
		WindowQuery{
			Namespace:   a.config.Namespace,
			Service:     a.config.Service,
			WindowStart: windowStart,
			WindowEnd:   windowEnd,
		},
	)
	if err != nil {
		return err
	}

	grouped := make(map[string][]*SketchRecord)

	for _, record := range records {
		if err := a.validateRecord(record); err != nil {
			a.logger.Warn(
				"[Go-Sail] <monitoring> ignore invalid distribution record",
				zap.String("recordID", record.RecordID),
				zap.Error(err),
			)
			continue
		}

		grouped[record.MetricKey] = append(
			grouped[record.MetricKey],
			record,
		)
	}

	for _, group := range grouped {
		if err := a.aggregateGroup(
			ctx,
			windowStart,
			windowEnd,
			group,
		); err != nil {
			return err
		}
	}

	return nil
}

func (a *Aggregator) aggregateGroup(ctx context.Context, windowStart time.Time, windowEnd time.Time, records []*SketchRecord) error {
	if len(records) == 0 {
		return nil
	}

	first := records[0]

	sketch, err := codec.NewDDSketch(first.RelativeAccuracy)
	if err != nil {
		return err
	}

	nodeSet := make(map[string]struct{})
	recordSet := make(map[string]struct{})

	var (
		sampleCount uint64
		sum         float64
		min         = math.Inf(1)
		max         = math.Inf(-1)
		recordIDs   = make([]string, 0, len(records))
	)

	for _, record := range records {
		// 防止 Store 异常返回重复记录。
		if _, exists := recordSet[record.RecordID]; exists {
			continue
		}
		recordSet[record.RecordID] = struct{}{}

		if err := sketch.MergeEncoded(record.Payload); err != nil {
			return fmt.Errorf(
				"[Go-Sail] <monitoring> merge record %s: %w",
				record.RecordID,
				err,
			)
		}

		nodeSet[record.NodeID] = struct{}{}
		recordIDs = append(recordIDs, record.RecordID)

		sampleCount += record.Count
		sum += record.Sum

		if record.Min < min {
			min = record.Min
		}
		if record.Max > max {
			max = record.Max
		}
	}

	if sampleCount == 0 {
		return nil
	}

	values, err := sketch.Quantiles(a.config.Quantiles)
	if err != nil {
		return err
	}

	quantiles := make(map[string]float64, len(values))
	for index, quantile := range a.config.Quantiles {
		quantiles[formatQuantile(quantile)] = values[index]
	}

	result := &AggregatedResult{
		SchemaVersion: SchemaVersionV1,
		Algorithm:     AlgorithmDDSketchV1,

		ResultID: BuildResultID(
			first.MetricKey,
			windowStart,
			windowEnd,
		),
		MetricKey: first.MetricKey,

		Namespace: first.Namespace,
		Service:   first.Service,
		Metric:    first.Metric,
		Labels:    first.Labels.Clone(),

		WindowStart: windowStart.UTC(),
		WindowEnd:   windowEnd.UTC(),

		RelativeAccuracy: first.RelativeAccuracy,

		NodeCount:   uint64(len(nodeSet)),
		RecordCount: uint64(len(recordSet)),
		SampleCount: sampleCount,

		Sum: sum,
		Min: min,
		Max: max,
		Avg: sum / float64(sampleCount),

		Quantiles: quantiles,

		// 保留聚合窗口级别的合并 Sketch，供下游继续合并。
		Payload:   sketch.Encode(),
		Finalized: true,
		CreatedAt: time.Now().UTC(),
	}

	if err := a.store.SaveResult(ctx, result); err != nil {
		return err
	}

	return a.store.MarkAggregated(ctx, recordIDs)
}

func (a *Aggregator) validateRecord(record *SketchRecord) error {
	switch {
	case record.SchemaVersion != SchemaVersionV1:
		return fmt.Errorf(
			"[Go-Sail] <monitoring> unsupported schema version: %d",
			record.SchemaVersion,
		)

	case record.Algorithm != AlgorithmDDSketchV1:
		return fmt.Errorf(
			"[Go-Sail] <monitoring> unsupported algorithm: %s",
			record.Algorithm,
		)

	case record.Namespace != a.config.Namespace:
		return fmt.Errorf(
			"[Go-Sail] <monitoring> unexpected namespace: %s",
			record.Namespace,
		)

	case record.Service != a.config.Service:
		return fmt.Errorf(
			"[Go-Sail] <monitoring> unexpected service: %s",
			record.Service,
		)

	case record.RelativeAccuracy != a.config.RelativeAccuracy:
		return fmt.Errorf(
			"[Go-Sail] <monitoring> relative accuracy mismatch: %v",
			record.RelativeAccuracy,
		)

	case record.Unit != UnitMicrosecond:
		return fmt.Errorf("[Go-Sail] <monitoring> unexpected unit: %s", record.Unit)

	case record.RecordID != record.BuildRecordID():
		return fmt.Errorf("[Go-Sail] <monitoring> invalid record ID")
	}

	return nil
}

func formatQuantile(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
