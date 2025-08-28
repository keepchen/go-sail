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

// RedisClientSetter redis连接实例设置器
type RedisClientSetter interface {
	// ForSchedule 为计划任务指定redis连接实例
	//
	// 调用此方法后，默认的连接实例将被覆盖
	ForSchedule()
	// ForRedisLocker 为redis分布式锁指定redis连接实例
	//
	// 调用此方法后，默认的连接实例将被覆盖
	ForRedisLocker()
}

// SetRedisClient 设置redis连接实例
func SetRedisClient(client redisLib.UniversalClient) RedisClientSetter {
	return &redisClientSetter{redisClient: client}
}

// ForSchedule 为计划任务指定redis连接实例
//
// 调用此方法后，默认的连接实例将被覆盖
func (s *redisClientSetter) ForSchedule() {
	schedule.SetRedisClientOnce(s.redisClient)
}

// ForRedisLocker 为redis分布式锁指定redis连接实例
//
// 调用此方法后，默认的连接实例将被覆盖
func (s *redisClientSetter) ForRedisLocker() {
	setRedisClientForRedisLockerOnce.Do(func() {
		redisClientForRedisLocker = s.redisClient
	})
}
