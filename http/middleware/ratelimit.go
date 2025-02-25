package middleware

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	redisLib "github.com/go-redis/redis/v8"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/api"
)

// Limiter 限流器支持本地和Redis两种模式
type Limiter struct {
	reqs   int           //每个时间窗口允许的最大请求数
	window time.Duration //滑动时间窗口长度

	//本地限流使用
	ips sync.Map

	//Redis客户端限流使用
	redisClient    redisLib.UniversalClient
	redisKeyPrefix string
}

// LimiterOptions 限流器设置项
type LimiterOptions struct {
	Reqs           int                      //请求数
	Window         time.Duration            //时间窗口
	RedisClient    redisLib.UniversalClient //Redis客户端实例，支持Client和ClusterClient
	RedisKeyPrefix string                   //Redis键名前缀
}

// NewLimiter 创建限流器，支持本地堆栈和Redis
func NewLimiter(opts LimiterOptions) *Limiter {
	limiter := &Limiter{
		reqs:           opts.Reqs,
		window:         opts.Window,
		redisClient:    opts.RedisClient,
		redisKeyPrefix: opts.RedisKeyPrefix,
	}

	return limiter
}

// AllowResult 检查结果
type AllowResult struct {
	Allowed   bool
	Remaining int
	ResetTime time.Time
}

// Allow 检查请求是否允许
func (l *Limiter) Allow(ip string) AllowResult {
	if l.redisClient != nil {
		return l.allowWithRedis(ip)
	}
	return l.allowWithLocal(ip)
}

// 本地限流逻辑
func (l *Limiter) allowWithLocal(ip string) AllowResult {
	now := time.Now()
	cutoff := now.Add(-l.window)
	var queue []time.Time
	if val, ok := l.ips.Load(ip); ok {
		queue = val.([]time.Time)
	}

	//移除过期请求
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

// Redis 限流逻辑脚本
var redisScript = redisLib.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local member = tonumber(ARGV[4])

-- 清除窗口外的请求
redis.call("ZREMRANGEBYSCORE", key, 0, now - window)

-- 计算当前请求数
local count = redis.call("ZCARD", key)

if count < limit then
    -- 添加当前请求时间
    redis.call("ZADD", key, now, member)
    -- 仅在新键创建时设置 TTL
    if redis.call("TTL", key) == -1 then
        redis.call("EXPIRE", key, window)  -- 使用秒
    end
    return {1, limit - (count + 1), now + window}
else
    -- 获取窗口重置时间
    local oldest = redis.call("ZRANGE", key, 0, 0, "WITHSCORES")
    local resetTime = oldest and oldest[2] and tonumber(oldest[2]) + window or now + window
    return {0, limit - count, resetTime}
end
`)

func (l *Limiter) allowWithRedis(ip string) AllowResult {
	ctx := context.Background()
	key := fmt.Sprintf("%s:{%s}", l.redisKeyPrefix, ip)
	now := time.Now().Unix()
	member := time.Now().UnixNano() + rand.Int64N(int64(l.reqs))

	result, err := redisScript.Run(ctx, l.redisClient, []string{key}, l.reqs, int(l.window.Seconds()), now, member).Result()
	if err != nil {
		//如果Redis操作失败，默认允许请求
		return AllowResult{Allowed: true, Remaining: l.reqs - 1, ResetTime: time.Now().Add(l.window)}
	}

	values, ok := result.([]interface{})
	if !ok || len(values) < 3 {
		//如果返回值格式不正确，默认拒绝请求
		return AllowResult{Allowed: false, Remaining: 0, ResetTime: time.Now().Add(l.window)}
	}

	allowed := values[0].(int64) == 1
	remaining := int(values[1].(int64))
	resetTime := time.Unix(values[2].(int64), 0)

	return AllowResult{allowed, remaining, resetTime}
}

// RateLimiter 中间件函数
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
