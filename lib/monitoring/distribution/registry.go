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

// Registry 管理本地指标记录器，并将已完成的 Sketch 放入有界非阻塞上报队列。
// Registry 支持并发使用。
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

// NewRegistry 应用默认值、校验配置并创建空注册表。
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

// ObserveDuration 以微秒记录耗时，并保留小数部分。
func (r *Registry) ObserveDuration(metric string, duration time.Duration, labels Labels) {
	// 内部统一使用微秒，但保留小数。
	value := float64(duration.Nanoseconds()) / float64(time.Microsecond)

	r.Observe(metric, value, labels)
}

// Observe 记录一个有限非负数值。无效观测和记录器错误不会返回给调用方，
// 而是由 ObserveErrors 计数。
func (r *Registry) Observe(metric string, value float64, labels Labels) {
	entry, err := r.getOrCreate(metric, labels)
	if err != nil {
		r.observeErrors.Add(1)
		return
	}

	// 除定时轮转外，观测路径也执行轮转，避免边界后到达的数值落入上一窗口。
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

// Start 持续轮转本地窗口，直至 ctx 被取消。关闭时会刷新当前不完整窗口，
// 将其放入队列后关闭 ReportQueue。
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

// ReportQueue 返回供 ReportWorker 消费的记录流。
// Start 会在关闭阶段关闭该通道。
func (r *Registry) ReportQueue() <-chan *SketchRecord {
	return r.reportQueue
}

// DroppedRecords 返回因上报无法立即接收或关闭超时而丢弃的记录数。
func (r *Registry) DroppedRecords() uint64 {
	return r.droppedRecords.Load()
}

// ObserveErrors 返回被记录器拒绝的观测次数。
func (r *Registry) ObserveErrors() uint64 {
	return r.observeErrors.Load()
}
