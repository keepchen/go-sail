package distribution

import (
	"context"
	"time"
)

// WindowQuery 查询完全包含在半开时间窗口内的节点级记录。
type WindowQuery struct {
	Namespace string // 待查询指标记录所属的命名空间。
	Service   string // 待查询指标记录所属的服务。

	WindowStart time.Time // 闭区间下界。
	WindowEnd   time.Time // 开区间上界。
}

// Store 持久化节点级 Sketch、聚合结果和补算游标。
type Store interface {
	Reporter

	// ListRecords 返回完全包含在查询窗口内的记录。
	ListRecords(ctx context.Context, query WindowQuery) ([]*SketchRecord, error)

	// SaveResult 按 ResultID 幂等保存聚合结果。
	SaveResult(ctx context.Context, result *AggregatedResult) error

	// GetResult 根据稳定标识读取聚合结果。
	GetResult(ctx context.Context, resultID string) (*AggregatedResult, error)

	// MarkAggregated 标记指定记录已参与聚合。
	MarkAggregated(ctx context.Context, recordIDs []string) error

	// GetAggregationCursor 返回最近成功处理窗口的结束时间。
	// 零值时间表示尚未持久化游标。
	GetAggregationCursor(ctx context.Context, namespace string, service string) (time.Time, error)

	// SaveAggregationCursor 推进已成功处理的窗口边界。
	SaveAggregationCursor(ctx context.Context, namespace string, service string, windowEnd time.Time) error
}
