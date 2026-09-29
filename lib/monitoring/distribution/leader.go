package distribution

import (
	"context"
	"time"
)

// Leadership 表示用于选举唯一聚合器的可续期租约。
type Leadership interface {
	// Acquire 尝试立即获取租约，不等待其他持有者释放。
	Acquire(ctx context.Context) (bool, error)
	// Renew 延长当前实例仍持有的租约。
	Renew(ctx context.Context) (bool, error)
	// Release 释放当前实例仍持有的租约。
	Release(ctx context.Context) error
}

// LeaderFactory 为存储后端创建选主租约。
type LeaderFactory interface {
	// NewLeadership 使用 key 和 ttl 创建租约。
	NewLeadership(key string, ttl time.Duration) Leadership
}
