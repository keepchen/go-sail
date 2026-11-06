package sail

import (
	"sync"

	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/schedule"
)

var (
	setRedisClientForRedisLockerOnce = &sync.Once{}
)

type redisClientSetter struct {
	redisClient redisLib.UniversalClient
}

// RedisClientSetter Redis connection instance setter
type RedisClientSetter interface {
	// ForSchedule specifies the Redis connection instance for scheduled tasks
	//
	// Calling this method overrides the default connection instance
	ForSchedule()
	// ForRedisLocker specifies the Redis connection instance for Redis distributed locks
	//
	// Calling this method overrides the default connection instance
	ForRedisLocker()
}

// SetRedisClient sets the Redis connection instance
func SetRedisClient(client redisLib.UniversalClient) RedisClientSetter {
	return &redisClientSetter{redisClient: client}
}

// ForSchedule specifies the Redis connection instance for scheduled tasks
//
// Calling this method overrides the default connection instance
func (s *redisClientSetter) ForSchedule() {
	schedule.SetRedisClientOnce(s.redisClient)
}

// ForRedisLocker specifies the Redis connection instance for Redis distributed locks
//
// Calling this method overrides the default connection instance
func (s *redisClientSetter) ForRedisLocker() {
	setRedisClientForRedisLockerOnce.Do(func() {
		redisClientForRedisLocker = s.redisClient
	})
}
