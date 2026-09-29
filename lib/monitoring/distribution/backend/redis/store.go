package redisbackend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	redisLib "github.com/go-redis/redis/v8"

	"github.com/keepchen/go-sail/v3/lib/monitoring/distribution"
)

// Store 是基于 Redis 的分布式统计存储。原始记录和聚合结果拥有独立保留期，
// 聚合游标不会过期。
type Store struct {
	client redisLib.UniversalClient

	localWindow     time.Duration
	recordRetention time.Duration
	resultRetention time.Duration
}

// NewStore 创建 Redis 存储。省略 resultRetention 或传入非正值时，
// 聚合结果使用 recordRetention，以保持向后兼容。
func NewStore(client redisLib.UniversalClient, localWindow time.Duration, recordRetention time.Duration, resultRetention ...time.Duration) *Store {
	resultTTL := recordRetention
	if len(resultRetention) > 0 && resultRetention[0] > 0 {
		resultTTL = resultRetention[0]
	}
	return &Store{
		client:          client,
		localWindow:     localWindow,
		recordRetention: recordRetention,
		resultRetention: resultTTL,
	}
}

// Report 幂等保存记录，并按本地窗口建立索引。
func (s *Store) Report(ctx context.Context, record *distribution.SketchRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> marshal sketch record: %w", err)
	}

	windowStart := record.WindowStart.Truncate(s.localWindow)
	indexKey := windowKey(record.Namespace, record.Service, windowStart)

	pipe := s.client.TxPipeline()

	// 稳定的 RecordID 使重复上报覆盖同一份记录。
	pipe.Set(ctx, recordKey(record.RecordID), data, s.recordRetention)
	pipe.SAdd(ctx, indexKey, record.RecordID)
	pipe.Expire(ctx, indexKey, s.recordRetention)

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> store sketch record: %w", err)
	}

	return nil
}

// ListRecords 读取完全包含在查询半开窗口内的记录。
func (s *Store) ListRecords(ctx context.Context, query distribution.WindowQuery) ([]*distribution.SketchRecord, error) {
	ids := make(map[string]struct{})

	for cursor := query.WindowStart; cursor.Before(query.WindowEnd); cursor = cursor.Add(s.localWindow) {
		key := windowKey(query.Namespace, query.Service, cursor)

		values, err := s.client.SMembers(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("[Go-Sail] <monitoring> list window records: %w", err)
		}

		for _, value := range values {
			ids[value] = struct{}{}
		}
	}

	if len(ids) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(ids))
	for id := range ids {
		keys = append(keys, recordKey(id))
	}

	values, err := s.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> load window records: %w", err)
	}

	records := make([]*distribution.SketchRecord, 0, len(values))

	for _, value := range values {
		if value == nil {
			continue
		}

		raw, ok := value.(string)
		if !ok {
			continue
		}

		var record distribution.SketchRecord
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			return nil, fmt.Errorf("[Go-Sail] <monitoring> unmarshal sketch record: %w", err)
		}

		if record.WindowStart.Before(query.WindowStart) ||
			!record.WindowEnd.Before(query.WindowEnd.Add(time.Nanosecond)) {
			continue
		}

		records = append(records, &record)
	}

	return records, nil
}

// SaveResult 使用结果保留期幂等保存聚合结果。
func (s *Store) SaveResult(ctx context.Context, result *distribution.AggregatedResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> marshal aggregate result: %w", err)
	}

	if err := s.client.Set(ctx, resultKey(result.ResultID), data, s.resultRetention).Err(); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> save aggregate result: %w", err)
	}

	return nil
}

// GetResult 根据 ID 读取聚合结果。
func (s *Store) GetResult(ctx context.Context, resultID string) (*distribution.AggregatedResult, error) {
	data, err := s.client.Get(ctx, resultKey(resultID)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> load aggregate result: %w", err)
	}

	var result distribution.AggregatedResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("[Go-Sail] <monitoring> unmarshal aggregate result: %w", err)
	}

	return &result, nil
}

// MarkAggregated 为已处理的记录 ID 创建带保留期的标记。
func (s *Store) MarkAggregated(ctx context.Context, recordIDs []string) error {
	if len(recordIDs) == 0 {
		return nil
	}

	pipe := s.client.Pipeline()

	for _, recordID := range recordIDs {
		pipe.Set(ctx, aggregatedKey(recordID), "1", s.recordRetention)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> mark records aggregated: %w", err)
	}

	return nil
}

// GetAggregationCursor 返回最近成功处理窗口的结束时间；
// 聚合尚未建立游标时返回零值时间。
func (s *Store) GetAggregationCursor(ctx context.Context, namespace string, service string) (time.Time, error) {
	value, err := s.client.Get(ctx, cursorKey(namespace, service)).Int64()
	if err != nil {
		if errors.Is(err, redisLib.Nil) {
			return time.Time{}, nil
		}
		return time.Time{}, fmt.Errorf("[Go-Sail] <monitoring> load aggregation cursor: %w", err)
	}
	return time.Unix(0, value).UTC(), nil
}

// SaveAggregationCursor 持久化最近成功处理窗口的结束时间。
func (s *Store) SaveAggregationCursor(ctx context.Context, namespace string, service string, windowEnd time.Time) error {
	if err := s.client.Set(ctx, cursorKey(namespace, service), windowEnd.UnixNano(), 0).Err(); err != nil {
		return fmt.Errorf("[Go-Sail] <monitoring> save aggregation cursor: %w", err)
	}
	return nil
}
