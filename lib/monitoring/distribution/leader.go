package distribution

import (
	"context"
	"time"
)

// Leadership represents a renewable lease for electing a single aggregator.
type Leadership interface {
	// Acquire attempts to acquire the lease immediately without waiting for other holders to release it.
	Acquire(ctx context.Context) (bool, error)
	// Renew extends the lease still held by the current instance.
	Renew(ctx context.Context) (bool, error)
	// Release releases the lease still held by the current instance.
	Release(ctx context.Context) error
}

// LeaderFactory creates leader election leases for a storage backend.
type LeaderFactory interface {
	// NewLeadership creates a lease using key and ttl.
	NewLeadership(key string, ttl time.Duration) Leadership
}
