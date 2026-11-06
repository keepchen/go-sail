package distribution

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

// Reporter persists or transmits a node-level Sketch record.
// Implementations should ensure idempotency based on SketchRecord.RecordID.
type Reporter interface {
	// Report writes a Sketch record.
	Report(ctx context.Context, record *SketchRecord) error
}

// ReportWorker consumes the registry queue and reports records with bounded retries.
type ReportWorker struct {
	reporter Reporter
	input    <-chan *SketchRecord
	logger   *zap.Logger

	maxAttempts uint
	retryDelay  time.Duration

	successes atomic.Uint64
	failures  atomic.Uint64
}

// NewReportWorker creates a reporting worker with the default retry policy.
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

// Run continuously consumes records until the input closes or ctx is canceled.
// After cancellation, it attempts to drain records already queued in Registry within a limited time.
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
			// Wait for the next retry.
		}
	}

	return errors.Join(
		errors.New("[Go-Sail] <monitoring> report attempts exhausted"),
		lastErr,
	)
}
