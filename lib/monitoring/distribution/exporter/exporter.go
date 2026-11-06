package exporter

import (
	"context"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
)

// Exporter publishes finalized aggregated results to external destinations.
type Exporter interface {
	// Export publishes an aggregated result.
	Export(ctx context.Context, result *distribution.AggregatedResult) error
}

// MultiExporter invokes multiple exporters sequentially in declaration order.
type MultiExporter struct {
	exporters []Exporter
}

var _ Exporter = (*MultiExporter)(nil)

// NewMultiExporter creates an exporter chain that runs sequentially.
func NewMultiExporter(exporters ...Exporter) *MultiExporter {
	return &MultiExporter{
		exporters: exporters,
	}
}

// Export calls all exporters in order and stops at the first error.
func (m *MultiExporter) Export(ctx context.Context, result *distribution.AggregatedResult) error {
	for _, item := range m.exporters {
		if err := item.Export(ctx, result); err != nil {
			return err
		}
	}
	return nil
}
