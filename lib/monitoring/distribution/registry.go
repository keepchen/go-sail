package distribution

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type metricEntry struct {
	metric    string
	labels    Labels
	metricKey string
	recorder  *Recorder
}

// Registry manages local metric recorders and enqueues completed Sketches in a bounded nonblocking reporting queue.
// Registry supports concurrent use.
type Registry struct {
	config *Conf

	mu      sync.RWMutex
	entries map[string]*metricEntry

	reportQueue  chan *SketchRecord
	closeOnce    sync.Once
	reportMu     sync.RWMutex
	reportClosed bool

	droppedRecords atomic.Uint64
	observeErrors  atomic.Uint64
}

// NewRegistry applies defaults, validates configuration, and creates an empty registry.
func NewRegistry(config *Conf) (*Registry, error) {
	config.SetDefaults()

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Registry{
		config:      config,
		entries:     make(map[string]*metricEntry),
		reportQueue: make(chan *SketchRecord, config.Reporter.QueueSize),
	}, nil
}

// ObserveDuration records durations in microseconds, preserving fractional values.
func (r *Registry) ObserveDuration(metric string, duration time.Duration, labels Labels) {
	// Use microseconds internally while preserving fractional values.
	value := float64(duration.Nanoseconds()) / float64(time.Microsecond)

	r.Observe(metric, value, labels)
}

// Observe records a finite nonnegative value. Invalid observations and recorder errors are not returned to the caller;
// instead, they are counted by ObserveErrors.
func (r *Registry) Observe(metric string, value float64, labels Labels) {
	entry, err := r.getOrCreate(metric, labels)
	if err != nil {
		r.observeErrors.Add(1)
		return
	}

	// In addition to scheduled rotation, the observation path rotates windows to prevent values arriving after a boundary from entering the previous window.
	r.rotateEntry(entry, time.Now())
	if err := entry.recorder.Observe(value); err != nil {
		r.observeErrors.Add(1)
	}
}

func (r *Registry) getOrCreate(metric string, labels Labels) (*metricEntry, error) {
	metricKey := BuildMetricKey(
		r.config.Namespace,
		r.config.Service,
		metric,
		labels,
	)

	r.mu.RLock()
	entry := r.entries[metricKey]
	r.mu.RUnlock()

	if entry != nil {
		return entry, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if entry = r.entries[metricKey]; entry != nil {
		return entry, nil
	}

	recorder, err := NewRecorder(
		time.Now(),
		r.config.LocalWindow,
		r.config.RelativeAccuracy,
	)
	if err != nil {
		return nil, err
	}

	entry = &metricEntry{
		metric:    metric,
		labels:    labels.Clone(),
		metricKey: metricKey,
		recorder:  recorder,
	}

	r.entries[metricKey] = entry

	return entry, nil
}

// Start continuously rotates local windows until ctx is canceled. On shutdown, it flushes the current incomplete window,
// enqueues it, and then closes ReportQueue.
func (r *Registry) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.flush()
			r.closeReportQueue()
			return

		case now := <-ticker.C:
			r.rotate(now)
		}
	}
}

func (r *Registry) rotate(now time.Time) {
	r.mu.RLock()
	entries := make([]*metricEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		entries = append(entries, entry)
	}
	r.mu.RUnlock()

	for _, entry := range entries {
		r.rotateEntry(entry, now)
	}
}

func (r *Registry) rotateEntry(entry *metricEntry, now time.Time) {
	snapshot, err := entry.recorder.Rotate(now)
	if err == nil && snapshot != nil {
		r.enqueue(r.buildRecord(entry, snapshot))
	}
}

func (r *Registry) flush() {
	r.mu.RLock()
	entries := make([]*metricEntry, 0, len(r.entries))
	for _, entry := range r.entries {
		entries = append(entries, entry)
	}
	r.mu.RUnlock()

	deadline := time.NewTimer(10 * time.Second)
	defer deadline.Stop()
	for _, entry := range entries {
		snapshot, err := entry.recorder.Flush()
		if err == nil && snapshot != nil {
			if !r.enqueueBefore(r.buildRecord(entry, snapshot), deadline.C) {
				return
			}
		}
	}
}

func (r *Registry) enqueue(record *SketchRecord) {
	r.reportMu.RLock()
	defer r.reportMu.RUnlock()
	if r.reportClosed {
		r.droppedRecords.Add(1)
		return
	}
	select {
	case r.reportQueue <- record:
	default:
		r.droppedRecords.Add(1)
	}
}

func (r *Registry) enqueueBefore(record *SketchRecord, deadline <-chan time.Time) bool {
	r.reportMu.RLock()
	defer r.reportMu.RUnlock()
	if r.reportClosed {
		r.droppedRecords.Add(1)
		return false
	}
	select {
	case r.reportQueue <- record:
		return true
	case <-deadline:
		r.droppedRecords.Add(1)
		return false
	}
}

func (r *Registry) closeReportQueue() {
	r.closeOnce.Do(func() {
		r.reportMu.Lock()
		defer r.reportMu.Unlock()
		r.reportClosed = true
		close(r.reportQueue)
	})
}

func (r *Registry) buildRecord(entry *metricEntry, snapshot *RecorderSnapshot) *SketchRecord {
	record := &SketchRecord{
		SchemaVersion: SchemaVersionV1,
		Algorithm:     AlgorithmDDSketchV1,

		MetricKey: entry.metricKey,
		Namespace: r.config.Namespace,
		Service:   r.config.Service,
		Metric:    entry.metric,
		Labels:    entry.labels.Clone(),
		NodeID:    r.config.NodeID,

		WindowStart: snapshot.WindowStart.UTC(),
		WindowEnd:   snapshot.WindowEnd.UTC(),

		Unit:             UnitMicrosecond,
		RelativeAccuracy: r.config.RelativeAccuracy,

		Count:   snapshot.Count,
		Sum:     snapshot.Sum,
		Min:     snapshot.Min,
		Max:     snapshot.Max,
		Payload: snapshot.Sketch.Encode(),
	}

	record.RecordID = record.BuildRecordID()

	return record
}

// ReportQueue returns the record stream consumed by ReportWorker.
// Start closes this channel during shutdown.
func (r *Registry) ReportQueue() <-chan *SketchRecord {
	return r.reportQueue
}

// DroppedRecords returns the number of records dropped because reporting could not accept them immediately or shutdown timed out.
func (r *Registry) DroppedRecords() uint64 {
	return r.droppedRecords.Load()
}

// ObserveErrors returns the number of observations rejected by recorders.
func (r *Registry) ObserveErrors() uint64 {
	return r.observeErrors.Load()
}
