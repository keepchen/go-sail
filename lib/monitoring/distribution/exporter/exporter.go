package exporter

import (
	"context"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
)

// Exporter 将已终态化的聚合结果发布到外部目标。
type Exporter interface {
	// Export 发布一条聚合结果。
	Export(ctx context.Context, result *distribution.AggregatedResult) error
}

// MultiExporter 按声明顺序依次调用多个导出器。
type MultiExporter struct {
	exporters []Exporter
}

var _ Exporter = (*MultiExporter)(nil)

// NewMultiExporter 创建顺序执行的导出器链。
func NewMultiExporter(exporters ...Exporter) *MultiExporter {
	return &MultiExporter{
		exporters: exporters,
	}
}

// Export 依次调用所有导出器，并在遇到首个错误时停止。
func (m *MultiExporter) Export(ctx context.Context, result *distribution.AggregatedResult) error {
	for _, item := range m.exporters {
		if err := item.Export(ctx, result); err != nil {
			return err
		}
	}
	return nil
}
