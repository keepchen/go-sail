package sail

import (
	"sync"

	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/utils"
)

var (
	redisClient redisLib.UniversalClient
	once        = &sync.Once{}
)

// SetRedisClientOnce 设置redis连接客户端
func SetRedisClientOnce(client redisLib.UniversalClient) {
	once.Do(func() {
		redisClient = client
	})
}

// RedisLocker 获取redis锁实例
func RedisLocker() utils.IRedisLocker {
	if redisClient != nil {
		return utils.RedisLocker(redisClient)
	}

	return utils.RedisLocker()
}
