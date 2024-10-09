package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/api"
)

// request 结构体用于存储每个请求的时间戳
type request struct {
	timestamp time.Time
}

// Limiter 是一个基于滑动时间窗口和IP限制的限流器
type Limiter struct {
	reqs   int                  // 每个时间窗口允许的最大请求数
	window time.Duration        // 滑动时间窗口的长度
	mu     sync.Mutex           // 互斥锁保护IP记录
	ips    map[string][]request // 每个IP地址的请求队列
}

// NewLimiter 返回一个新的基于滑动时间窗口和IP限制的Limiter实例
func NewLimiter(reqs int, window time.Duration) *Limiter {
	return &Limiter{
		reqs:   reqs,
		window: window,
		ips:    make(map[string][]request),
	}
}

// Allow 方法返回指定IP的请求是否被允许，同时返回剩余的请求数和窗口重置时间
func (l *Limiter) Allow(ip string) (bool, int, time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	//如果IP不存在，则初始化
	if _, exists := l.ips[ip]; !exists {
		l.ips[ip] = []request{}
	}

	//移除过期的请求
	cutoff := now.Add(-l.window)
	queue := l.ips[ip]
	for len(queue) > 0 && queue[0].timestamp.Before(cutoff) {
		queue = queue[1:]
	}
	l.ips[ip] = queue

	//计算窗口重置时间
	var resetTime time.Time
	if len(l.ips[ip]) > 0 {
		resetTime = l.ips[ip][0].timestamp.Add(l.window)
	} else {
		resetTime = now.Add(l.window)
	}

	remaining := l.reqs - len(l.ips[ip])
	if len(l.ips[ip]) < l.reqs {
		//允许请求并记录请求时间
		l.ips[ip] = append(l.ips[ip], request{timestamp: now})
		return true, remaining - 1, resetTime
	}

	//否则拒绝请求
	return false, remaining, resetTime
}

// RateLimiter 是一个限流器中间件，用于对HTTP请求进行IP限流，并添加X-Rate-Limit相关的响应头
func RateLimiter(limiter *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		allowed, remaining, resetTime := limiter.Allow(ip)

		//设置X-Rate-Limit响应头
		c.Writer.Header().Set("X-Rate-Limit-Limit", strconv.Itoa(limiter.reqs))
		c.Writer.Header().Set("X-Rate-Limit-Remaining", strconv.Itoa(remaining))
		c.Writer.Header().Set("X-Rate-Limit-Reset", strconv.FormatInt(resetTime.Unix(), 10))

		if !allowed {
			api.Response(c).Wrap(constants.ErrRequestParamsInvalid, nil, "Too Many Request").
				SendWithCode(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
