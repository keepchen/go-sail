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

// Leadership 实现带所有权令牌的 Redis 租约。续期和释放使用先比较后执行的脚本，
// 防止一个进程修改另一个进程持有的租约。
type Leadership struct {
	client redisLib.UniversalClient

	key   string
	token string
	ttl   time.Duration
}

// NewLeadership 创建一个带随机所有权令牌的 Redis 租约。
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

// Acquire 在租约键不存在时尝试创建该键。
func (l *Leadership) Acquire(ctx context.Context) (bool, error) {
	acquired, err := l.client.SetNX(ctx, l.key, l.token, l.ttl).Result()

	if err != nil {
		return false, fmt.Errorf("[Go-Sail] <monitoring> acquire aggregator leadership: %w", err)
	}

	return acquired, nil
}

// Renew 仅在所有权令牌仍匹配时延长租约。
func (l *Leadership) Renew(ctx context.Context) (bool, error) {
	value, err := renewScript.Run(ctx, l.client, []string{l.key}, l.token, l.ttl.Milliseconds()).Int64()

	if err != nil {
		return false, fmt.Errorf("[Go-Sail] <monitoring> renew aggregator leadership: %w", err)
	}

	return value == 1, nil
}

// Release 仅在所有权令牌仍匹配时删除租约。
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
