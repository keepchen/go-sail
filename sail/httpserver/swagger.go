package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/sail/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RunSwaggerServerOnDebugMode 启动swagger文档服务
//
// 当配置文件指明启用时才会启动
func RunSwaggerServerOnDebugMode(conf config.SwaggerConf, ginEngine *gin.Engine) {
	if !conf.Enable {
		//如果不是调试模式就不注册swagger路由
		return
	}

	//swagger-ui
	ginEngine.StaticFile("/swagger-assets/doc.json", conf.JsonPath)
	url := ginSwagger.URL("/swagger-assets/doc.json") // The url pointing to API definition
	//access /swagger/index.html
	ginEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//redoc-ui
	ginEngine.StaticFile("/redoc/docs.html", conf.RedocUIPath)

	//favicon
	ginEngine.StaticFile("/favicon.ico", conf.FaviconPath)
}
