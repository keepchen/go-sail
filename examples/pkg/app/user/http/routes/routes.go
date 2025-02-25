package routes

import (
	"net/http"

	"github.com/gin-contrib/gzip"

	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/middleware"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/http/handler"
	mdlw "github.com/keepchen/go-sail/v3/http/middleware"
)

// RegisterRoutes 注册路由
func RegisterRoutes(r *gin.Engine) {
	//k8s健康检查接口
	r.GET("/actuator/health", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "ok")
	})
	allowHeaders := map[string]string{
		"Access-Control-Allow-Headers": "Authorization, Content-Type, Content-Length, Some-Other-Headers",
	}
	//limiter := mdlw.NewLimiter(mdlw.LimiterOptions{
	//	Reqs: 1000,
	//	Window: time.Minute,
	//	RedisClient: sail.GetRedis(),
	//	RedisKeyPrefix: "rate-limiter",
	//})
	//mdlw.RateLimiter(limiter)
	//全局打印请求载荷、放行跨域请求、gzip压缩
	r.Use(mdlw.LogTrace(), mdlw.PrintRequestPayload(), mdlw.WithCors(allowHeaders), gzip.Gzip(gzip.DefaultCompression))
	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/say-hello", handler.SayHello)
		userGroup := apiGroup.Group("/user")
		{
			userGroup.Use(middleware.AuthCheck()).GET("/info", handler.GetUserInfo)
		}
	}
}
