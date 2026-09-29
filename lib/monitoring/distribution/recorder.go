package distribution

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution/codec"
)

// RecorderSnapshot 是一个记录窗口的不可变交接快照。
// 快照返回后，原 Recorder 不会再修改其中的 Sketch。
type RecorderSnapshot struct {
	Sketch *codec.DDSketch // 观测分布的可合并近似摘要。

	WindowStart time.Time // 窗口闭区间边界。
	WindowEnd   time.Time // 窗口开区间边界。

	Count uint64  // 精确样本数。
	Sum   float64 // 精确样本总和。
	Min   float64 // 精确最小样本值。
	Max   float64 // 精确最大样本值。
}

// Recorder 将有限非负数值收集到对齐的本地窗口中，并支持并发使用。
type Recorder struct {
	mu sync.Mutex

	relativeAccuracy float64
	windowDuration   time.Duration

	sketch *codec.DDSketch

	windowStart time.Time
	windowEnd   time.Time

	count uint64
	sum   float64
	min   float64
	max   float64
}

// NewRecorder 创建一个与 now 所在时间窗口对齐的记录器。
func NewRecorder(now time.Time, windowDuration time.Duration, relativeAccuracy float64) (*Recorder, error) {
	sketch, err := codec.NewDDSketch(relativeAccuracy)
	if err != nil {
		return nil, err
	}

	windowStart := now.Truncate(windowDuration)

	return &Recorder{
		relativeAccuracy: relativeAccuracy,
		windowDuration:   windowDuration,
		sketch:           sketch,
		windowStart:      windowStart,
		windowEnd:        windowStart.Add(windowDuration),
		min:              math.Inf(1),
		max:              math.Inf(-1),
	}, nil
}

// Observe 向当前窗口添加一个有限非负数值。
func (r *Recorder) Observe(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return errors.New("[Go-Sail] <monitoring> observed value must be a finite non-negative number")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.sketch.Add(value); err != nil {
		return err
	}

	r.count++
	r.sum += value

	if value < r.min {
		r.min = value
	}
	if value > r.max {
		r.max = value
	}

	return nil
}

// Rotate 在 now 到达当前窗口末尾时关闭并返回该窗口。
// 窗口仍开放或不含样本时返回 (nil, nil)。
func (r *Recorder) Rotate(now time.Time) (*RecorderSnapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if now.Before(r.windowEnd) {
		return nil, nil
	}

	next, err := codec.NewDDSketch(r.relativeAccuracy)
	if err != nil {
		return nil, err
	}

	snapshot := &RecorderSnapshot{
		Sketch:      r.sketch,
		WindowStart: r.windowStart,
		WindowEnd:   r.windowEnd,
		Count:       r.count,
		Sum:         r.sum,
		Min:         r.min,
		Max:         r.max,
	}

	nextStart := now.Truncate(r.windowDuration)

	r.sketch = next
	r.windowStart = nextStart
	r.windowEnd = nextStart.Add(r.windowDuration)
	r.count = 0
	r.sum = 0
	r.min = math.Inf(1)
	r.max = math.Inf(-1)

	if snapshot.Count == 0 {
		return nil, nil
	}

	return snapshot, nil
}

// Flush 返回当前不完整窗口并重置记录器。快照保留原定窗口结束时间，
// 以确保分布式聚合仍按统一边界对齐。
func (r *Recorder) Flush() (*RecorderSnapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.count == 0 {
		return nil, nil
	}

	next, err := codec.NewDDSketch(r.relativeAccuracy)
	if err != nil {
		return nil, err
	}

	snapshot := &RecorderSnapshot{
		Sketch:      r.sketch,
		WindowStart: r.windowStart,
		WindowEnd:   r.windowEnd,
		Count:       r.count,
		Sum:         r.sum,
		Min:         r.min,
		Max:         r.max,
	}
	r.sketch = next
	r.count = 0
	r.sum = 0
	r.min = math.Inf(1)
	r.max = math.Inf(-1)

	return snapshot, nil
}
