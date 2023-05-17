package routes

import (
	"context"
	"net/http"

	"github.com/keepchen/go-sail/pkg/app/user/http/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/keepchen/go-sail/pkg/app/user/config"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/pkg/app/user/http/handler"
	mdlw "github.com/keepchen/go-sail/pkg/common/http/middleware"
)

// 注册路由
func registerRoutes(_ context.Context, r *gin.Engine) {
	if config.GetGlobalConfig().Debug {
		//仅在调试模式下才开始pprof检测
		pprof.Register(r, "debug/pprof")
	}

	//k8s健康检查接口
	r.GET("/actuator/health", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "ok")
	})

	allowHeaders := map[string]string{
		"Access-Control-Allow-Headers": "Authorization, Content-Type, Content-Length, Some-Other-Headers",
	}
	//全局打印请求载荷、放行跨域请求、写入Prometheus exporter
	r.Use(mdlw.Before(), mdlw.PrintRequestPayload(), mdlw.WithCors(allowHeaders), mdlw.PrometheusExporter())

	api := r.Group("/api/v1")
	{
		api.GET("/say-hello", handler.SayHello)

		userGroup := api.Group("/user")
		{
			userGroup.Use(middleware.AuthCheck()).GET("/info", handler.GetUserInfo)
		}
	}
}
