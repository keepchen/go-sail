package distribution

import (
	"context"
	"testing"
	"time"

	"go.uber.org/zap"
)

type cursorStore struct {
	cursor  time.Time
	queries []WindowQuery
}

func (s *cursorStore) Report(context.Context, *SketchRecord) error { return nil }
func (s *cursorStore) ListRecords(_ context.Context, query WindowQuery) ([]*SketchRecord, error) {
	s.queries = append(s.queries, query)
	return nil, nil
}
func (s *cursorStore) SaveResult(context.Context, *AggregatedResult) error { return nil }
func (s *cursorStore) GetResult(context.Context, string) (*AggregatedResult, error) {
	return nil, nil
}
func (s *cursorStore) MarkAggregated(context.Context, []string) error { return nil }
func (s *cursorStore) GetAggregationCursor(context.Context, string, string) (time.Time, error) {
	return s.cursor, nil
}
func (s *cursorStore) SaveAggregationCursor(_ context.Context, _, _ string, windowEnd time.Time) error {
	s.cursor = windowEnd
	return nil
}

func TestAggregatorBackfillsFromCursor(t *testing.T) {
	store := &cursorStore{cursor: time.Date(2026, 8, 12, 10, 1, 0, 0, time.UTC)}
	aggregator := NewAggregator(&Conf{
		Namespace:       "production",
		Service:         "api",
		AggregateWindow: time.Minute,
		AllowedLateness: 20 * time.Second,
	}, store, nil, zap.NewNop())

	now := time.Date(2026, 8, 12, 10, 5, 20, 0, time.UTC)
	if err := aggregator.aggregateClosedWindow(context.Background(), now); err != nil {
		t.Fatal(err)
	}

	if got, want := len(store.queries), 4; got != want {
		t.Fatalf("query count = %d, want %d", got, want)
	}
	if got, want := store.queries[0].WindowStart, store.cursor.Add(-4*time.Minute); !got.Equal(want) {
		t.Fatalf("first window start = %v, want %v", got, want)
	}
	if got, want := store.cursor, time.Date(2026, 8, 12, 10, 5, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("cursor = %v, want %v", got, want)
	}
}
