package middleware

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/api"
)

// Limiter supports local and Redis rate limiting
type Limiter struct {
	reqs   int           //Maximum requests allowed per time window
	window time.Duration //Sliding time window length

	//Used for local rate limiting
	ips sync.Map

	//Redis client for rate limiting
	redisClient    redisLib.UniversalClient
	redisKeyPrefix string
}

// LimiterOptions rate limiter options
type LimiterOptions struct {
	Reqs           int                      //Request count
	Window         time.Duration            //Time window
	RedisClient    redisLib.UniversalClient //Redis client instance; supports Client and ClusterClient
	RedisKeyPrefix string                   //Redis key prefix
}

// NewLimiter creates a rate limiter with local storage or Redis
func NewLimiter(opts LimiterOptions) *Limiter {
	limiter := &Limiter{
		reqs:           opts.Reqs,
		window:         opts.Window,
		redisClient:    opts.RedisClient,
		redisKeyPrefix: opts.RedisKeyPrefix,
	}

	return limiter
}

// AllowResult check result
type AllowResult struct {
	Allowed   bool
	Remaining int
	ResetTime time.Time
}

// Allow checks whether a request is allowed
func (l *Limiter) Allow(ip string) AllowResult {
	if l.redisClient != nil {
		return l.allowWithRedis(ip)
	}
	return l.allowWithLocal(ip)
}

// Local rate limiting logic
func (l *Limiter) allowWithLocal(ip string) AllowResult {
	now := time.Now()
	cutoff := now.Add(-l.window)
	var queue []time.Time
	if val, ok := l.ips.Load(ip); ok {
		queue = val.([]time.Time)
	}

	//Remove expired requests
	for len(queue) > 0 && queue[0].Before(cutoff) {
		queue = queue[1:]
	}

	remaining := l.reqs - len(queue)
	resetTime := now.Add(l.window)

	if remaining > 0 {
		queue = append(queue, now)
		l.ips.Store(ip, queue)
		return AllowResult{true, remaining - 1, resetTime}
	}
	return AllowResult{false, remaining, resetTime}
}

// Redis rate limiting script
var redisScript = redisLib.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local member = tonumber(ARGV[4])

-- Remove requests outside the window
redis.call("ZREMRANGEBYSCORE", key, 0, now - window)

-- Count current requests
local count = redis.call("ZCARD", key)

if count < limit then
    -- Add the current request timestamp
    redis.call("ZADD", key, now, member)
    -- Set TTL only when a new key is created
    if redis.call("TTL", key) == -1 then
        redis.call("EXPIRE", key, window)  -- Use seconds
    end
    return {1, limit - (count + 1), now + window}
else
    -- Get the window reset time
    local oldest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")
    local resetTime = oldest and oldest[2] and tonumber(oldest[2]) + window or now + window
    return {0, limit - count, resetTime}
end
`)

func (l *Limiter) allowWithRedis(ip string) AllowResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	key := fmt.Sprintf("%s:{%s}", l.redisKeyPrefix, ip)
	now := time.Now().Unix()
	member := time.Now().UnixNano() + int64(rand.Intn(l.reqs))

	result, err := redisScript.Run(ctx, l.redisClient, []string{key}, l.reqs, int(l.window.Seconds()), now, member).Result()
	if err != nil {
		//Allow the request by default if the Redis operation fails
		return AllowResult{Allowed: true, Remaining: l.reqs - 1, ResetTime: time.Now().Add(l.window)}
	}

	values, ok := result.([]any)
	if !ok || len(values) < 3 {
		//Reject the request by default if the return value has an invalid format
		return AllowResult{Allowed: false, Remaining: 0, ResetTime: time.Now().Add(l.window)}
	}

	allowed := values[0].(int64) == 1
	remaining := int(values[1].(int64))
	resetTime := time.Unix(values[2].(int64), 0)

	return AllowResult{allowed, remaining, resetTime}
}

// RateLimiter middleware function
func RateLimiter(limiter *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		res := limiter.Allow(ip)

		c.Writer.Header().Set("X-Rate-Limit-Limit", strconv.Itoa(limiter.reqs))
		c.Writer.Header().Set("X-Rate-Limit-Remaining", strconv.Itoa(res.Remaining))
		c.Writer.Header().Set("X-Rate-Limit-Reset", strconv.FormatInt(res.ResetTime.Unix(), 10))

		if !res.Allowed {
			api.Response(c).Wrap(constants.ErrRequestParamsInvalid, nil, "Too Many Request").
				SendWithCode(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
