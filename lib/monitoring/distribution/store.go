package distribution

import (
	"context"
	"time"
)

// WindowQuery queries node-level records fully contained in a half-open time window.
type WindowQuery struct {
	Namespace string // Namespace of the metric records to query.
	Service   string // Service of the metric records to query.

	WindowStart time.Time // Inclusive lower bound.
	WindowEnd   time.Time // Exclusive upper bound.
}

// Store persists node-level Sketches, aggregated results, and the backfill cursor.
type Store interface {
	Reporter

	// ListRecords returns records fully contained in the query window.
	ListRecords(ctx context.Context, query WindowQuery) ([]*SketchRecord, error)

	// SaveResult saves aggregated results idempotently by ResultID.
	SaveResult(ctx context.Context, result *AggregatedResult) error

	// GetResult reads an aggregated result by its stable identifier.
	GetResult(ctx context.Context, resultID string) (*AggregatedResult, error)

	// MarkAggregated marks the specified records as included in aggregation.
	MarkAggregated(ctx context.Context, recordIDs []string) error

	// GetAggregationCursor returns the end time of the most recently processed window.
	// The zero time indicates that no cursor has been persisted yet.
	GetAggregationCursor(ctx context.Context, namespace string, service string) (time.Time, error)

	// SaveAggregationCursor advances the boundary of successfully processed windows.
	SaveAggregationCursor(ctx context.Context, namespace string, service string, windowEnd time.Time) error
}
