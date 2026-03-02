package sail

import (
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/utils"
)

var redisClientForRedisLocker redisLib.UniversalClient

// RedisLocker 获取redis锁实例
//
// # Note
//
// 若没有初始化redis实例且没有设定redis实例，调用此方法将panic
//
// 调用者可以使用 SetRedisClient(client).ForRedisLocker() 覆盖默认连接
func RedisLocker() utils.IRedisLocker {
	if redisClientForRedisLocker != nil {
		return utils.RedisLocker(redisClientForRedisLocker)
	}

	return utils.RedisLocker()
}
