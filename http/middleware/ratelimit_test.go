package middleware

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/keepchen/go-sail/v3/lib/redis"
)

func TestNewLimiter(t *testing.T) {
	var (
		limit      = 100
		challenges = 500
	)
	t.Run("NewLimiter (with local)", func(t *testing.T) {
		limiter := NewLimiter(LimiterOptions{
			Reqs:   limit,
			Window: time.Minute,
		})
		for i := 0; i < challenges; i++ {
			result := limiter.Allow("127.0.0.1")
			t.Log(i, result.Allowed, result.Remaining, result.ResetTime)
			assert.Equal(t, i < limit, result.Allowed)
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
			Reqs:           limit,
			Window:         time.Minute,
			RedisClient:    redisClient,
			RedisKeyPrefix: fmt.Sprintf("%s-%d", "go-sail-rate-limiter", time.Now().UnixNano()),
		})
		for i := 0; i < challenges; i++ {
			result := limiter.Allow("127.0.0.1")
			t.Log(i, result.Allowed, result.Remaining, result.ResetTime)
			assert.Equal(t, i < limit, result.Allowed)

			result2 := limiter.Allow("192.168.100.1")
			t.Log(i, result2.Allowed, result2.Remaining, result2.ResetTime)
			assert.Equal(t, i < 10, result2.Allowed)
		}
	})

	t.Run("RateLimiter", func(t *testing.T) {
		c, _ := createTestContextAndEngine()

		req, _ := http.NewRequest("GET", "/test?name=foo", nil)
		req.Header.Set("Content-Type", "application/json")

		c.Request = req

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
			Reqs:           limit,
			Window:         time.Minute,
			RedisClient:    redisClient,
			RedisKeyPrefix: fmt.Sprintf("%s-%d", "go-sail-rate-limiter", time.Now().UnixNano()),
		})
		for i := 0; i < challenges; i++ {
			result := limiter.Allow("127.0.0.1")
			t.Log(i, result.Allowed, result.Remaining, result.ResetTime)
			assert.Equal(t, i < limit, result.Allowed)

			result2 := limiter.Allow("192.168.100.1")
			t.Log(i, result2.Allowed, result2.Remaining, result2.ResetTime)
			assert.Equal(t, i < limit, result2.Allowed)
		}

		RateLimiter(limiter)(c)
	})
}
