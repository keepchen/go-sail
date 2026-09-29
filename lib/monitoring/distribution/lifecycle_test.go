package distribution

import (
	"context"
	"testing"
	"time"
)

func TestRecorderFlushesPartialWindow(t *testing.T) {
	recorder, err := NewRecorder(time.Now(), time.Minute, 0.01)
	if err != nil {
		t.Fatal(err)
	}
	if err := recorder.Observe(42); err != nil {
		t.Fatal(err)
	}

	snapshot, err := recorder.Flush()
	if err != nil {
		t.Fatal(err)
	}
	if snapshot == nil || snapshot.Count != 1 || snapshot.Sum != 42 {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	if second, err := recorder.Flush(); err != nil || second != nil {
		t.Fatalf("second flush = %#v, %v; want nil, nil", second, err)
	}
}

func TestRegistryClosesWithFinalSnapshot(t *testing.T) {
	registry, err := NewRegistry(&Conf{
		Enabled:   true,
		Namespace: "production",
		Service:   "api",
		NodeID:    "node-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	registry.Observe("latency", 42, nil)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		registry.Start(ctx)
		close(done)
	}()
	cancel()
	<-done

	record, ok := <-registry.ReportQueue()
	if !ok || record.Count != 1 || record.Sum != 42 {
		t.Fatalf("unexpected final record: %#v, open=%v", record, ok)
	}
	if _, ok := <-registry.ReportQueue(); ok {
		t.Fatal("report queue was not closed")
	}
}
