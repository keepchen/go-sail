package redisbackend

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	redisLib "github.com/go-redis/redis/v8"
)

var renewScript = redisLib.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
return 0
`)

var releaseScript = redisLib.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
	return redis.call("DEL", KEYS[1])
end
return 0
`)

// Leadership implements a Redis lease with an ownership token. Renewal and release use compare-and-act scripts
// to prevent one process from modifying a lease held by another.
type Leadership struct {
	client redisLib.UniversalClient

	key   string
	token string
	ttl   time.Duration
}

// NewLeadership creates a Redis lease with a random ownership token.
func NewLeadership(client redisLib.UniversalClient, key string, ttl time.Duration) (*Leadership, error) {
	token, err := randomToken()
	if err != nil {
		return nil, err
	}

	return &Leadership{
		client: client,
		key:    key,
		token:  token,
		ttl:    ttl,
	}, nil
}

// Acquire attempts to create the lease key if it does not exist.
func (l *Leadership) Acquire(ctx context.Context) (bool, error) {
	acquired, err := l.client.SetNX(ctx, l.key, l.token, l.ttl).Result()

	if err != nil {
		return false, fmt.Errorf("[Go-Sail] <monitoring> acquire aggregator leadership: %w", err)
	}

	return acquired, nil
}

// Renew extends the lease only if the ownership token still matches.
func (l *Leadership) Renew(ctx context.Context) (bool, error) {
	value, err := renewScript.Run(ctx, l.client, []string{l.key}, l.token, l.ttl.Milliseconds()).Int64()

	if err != nil {
		return false, fmt.Errorf("[Go-Sail] <monitoring> renew aggregator leadership: %w", err)
	}

	return value == 1, nil
}

// Release deletes the lease only if the ownership token still matches.
func (l *Leadership) Release(ctx context.Context) error {
	_, err := releaseScript.Run(ctx, l.client, []string{l.key}, l.token).Result()

	if err != nil && !errors.Is(err, redisLib.Nil) {
		return fmt.Errorf("[Go-Sail] <monitoring> release aggregator leadership: %w", err)
	}

	return nil
}

func randomToken() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", fmt.Errorf("[Go-Sail] <monitoring> generate leader token: %w", err)
	}
	return hex.EncodeToString(value[:]), nil
}
