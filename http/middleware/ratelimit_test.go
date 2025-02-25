package middleware

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

func TestNewLimiter(t *testing.T) {
	t.Run("NewLimiter (with local)", func(t *testing.T) {
		limiter := NewLimiter(LimiterOptions{
			Reqs:   10,
			Window: time.Minute,
		})
		for i := 0; i < 15; i++ {
			result := limiter.Allow("127.0.0.1")
			t.Log(i, result.Allowed, result.Remaining, result.ResetTime)
			assert.Equal(t, i < 10, result.Allowed)
		}
	})

	t.Run("NewLimiter (with redis)", func(t *testing.T) {
		conf := redis.Conf{
			Endpoint: redis.Endpoint{
				Host: "127.0.0.1",
				Port: 6379,
			},
		}
		redisClient, err := redis.New(conf)
		if err != nil || redisClient == nil {
			t.Log("redis instance not ready, this test case ignore")
			return
		}
		limiter := NewLimiter(LimiterOptions{
			Reqs:           10,
			Window:         time.Minute,
			RedisClient:    redisClient,
			RedisKeyPrefix: fmt.Sprintf("%s-%d", "go-sail-rate-limiter", time.Now().UnixNano()),
		})
		for i := 0; i < 15; i++ {
			result := limiter.Allow("127.0.0.1")
			t.Log(i, result.Allowed, result.Remaining, result.ResetTime)
			assert.Equal(t, i < 10, result.Allowed)

			result2 := limiter.Allow("192.168.100.1")
			t.Log(i, result2.Allowed, result2.Remaining, result2.ResetTime)
			assert.Equal(t, i < 10, result2.Allowed)
		}
	})
}
