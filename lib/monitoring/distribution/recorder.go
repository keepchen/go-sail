package distribution

import (
	"errors"
	"math"
	"sync"
	"time"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution/codec"
)

// RecorderSnapshot is an immutable handoff snapshot of a recording window.
// The original Recorder will no longer modify the Sketch after returning the snapshot.
type RecorderSnapshot struct {
	Sketch *codec.DDSketch // Mergeable approximate summary of the observed distribution.

	WindowStart time.Time // Inclusive window boundary.
	WindowEnd   time.Time // Exclusive window boundary.

	Count uint64  // Exact sample count.
	Sum   float64 // Exact sample sum.
	Min   float64 // Exact minimum sample value.
	Max   float64 // Exact maximum sample value.
}

// Recorder collects finite nonnegative values into aligned local windows and supports concurrent use.
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

// NewRecorder creates a recorder aligned to the time window containing now.
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

// Observe adds a finite nonnegative value to the current window.
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

// Rotate closes and returns the current window when now reaches its end.
// Returns (nil, nil) if the window is still open or contains no samples.
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

// Flush returns the current incomplete window and resets the recorder. The snapshot retains the scheduled window end time
// to keep distributed aggregation aligned to common boundaries.
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
