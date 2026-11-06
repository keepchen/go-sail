package sail

import (
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/utils"
)

var redisClientForRedisLocker redisLib.UniversalClient

// RedisLocker retrieves the Redis lock instance
//
// # Note
//
// This method panics if no Redis instance has been initialized or explicitly set
//
// Callers can override the default connection with SetRedisClient(client).ForRedisLocker()
func RedisLocker() utils.IRedisLocker {
	if redisClientForRedisLocker != nil {
		return utils.RedisLocker(redisClientForRedisLocker)
	}

	return utils.RedisLocker()
}
