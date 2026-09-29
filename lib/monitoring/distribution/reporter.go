package distribution

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Reporter 持久化或传输一条节点级 Sketch 记录。
// 实现应针对 SketchRecord.RecordID 保证幂等。
type Reporter interface {
	// Report 写入一条 Sketch 记录。
	Report(ctx context.Context, record *SketchRecord) error
}

// ReportWorker 消费注册表队列，并以有限次数重试上报记录。
type ReportWorker struct {
	reporter Reporter
	input    <-chan *SketchRecord
	logger   *zap.Logger

	maxAttempts uint
	retryDelay  time.Duration

	successes atomic.Uint64
	failures  atomic.Uint64
}

// NewReportWorker 使用默认重试策略创建上报工作器。
func NewReportWorker(reporter Reporter, input <-chan *SketchRecord, logger *zap.Logger) *ReportWorker {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &ReportWorker{
		reporter:    reporter,
		input:       input,
		logger:      logger,
		maxAttempts: 3,
		retryDelay:  time.Second,
	}
}

// Run 持续消费记录，直至输入关闭或 ctx 被取消。
// 取消后会在限定时间内尝试排空 Registry 中已入队的记录。
func (w *ReportWorker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			w.drain()
			return

		case record, ok := <-w.input:
			if !ok {
				return
			}

			if err := w.report(ctx, record); err != nil {
				w.failures.Add(1)
				w.logger.Error(
					"[Go-Sail] <monitoring> report distribution sketch failed",
					zap.String("recordID", record.RecordID),
					zap.String("metric", record.Metric),
					zap.Error(err),
				)
				continue
			}

			w.successes.Add(1)
		}
	}
}

func (w *ReportWorker) drain() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		case record, ok := <-w.input:
			if !ok {
				return
			}
			if err := w.report(ctx, record); err != nil {
				w.failures.Add(1)
				w.logger.Error(
					"[Go-Sail] <monitoring> report distribution sketch during shutdown failed",
					zap.String("recordID", record.RecordID),
					zap.Error(err),
				)
			} else {
				w.successes.Add(1)
			}
		}
	}
}

func (w *ReportWorker) report(ctx context.Context, record *SketchRecord) error {
	var lastErr error

	for attempt := uint(1); attempt <= w.maxAttempts; attempt++ {
		if err := w.reporter.Report(ctx, record); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if attempt == w.maxAttempts {
			break
		}

		timer := time.NewTimer(w.retryDelay * time.Duration(attempt))

		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()

		case <-timer.C:
			// 等待下一次重试。
		}
	}

	return errors.Join(
		errors.New("[Go-Sail] <monitoring> report attempts exhausted"),
		lastErr,
	)
}
