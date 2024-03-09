package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithCorsOnlyOptions 允许浏览器跨域请求
//
// 仅对options探测请求注入放行headers
func WithCorsOnlyOptions(headers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var defaultCorsHeaders = map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Authorization, Content-Type, Content-Length",
			"Access-Control-Allow-Methods": "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD",
			"Access-Control-Expose-Headers": "Content-Length, Access-Control-Allow-Origin, " +
				"Access-Control-Allow-Headers, Content-Type, Authorization",
			"Access-Control-Allow-Credentials": "true",
		}
		if headers == nil {
			origin := c.Request.Header.Get("Origin")
			if len(origin) == 0 {
				origin = "*"
			}
			defaultCorsHeaders["Access-Control-Allow-Origin"] = origin
			headers = defaultCorsHeaders
		}
		for key, value := range headers {
			c.Writer.Header().Set(key, value)
		}

		//处理浏览器跨域options探测请求
		if c.Request.Method == http.MethodOptions {
			c.Writer.WriteHeader(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}

// WithCors 允许浏览器跨域请求
func WithCors(headers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var defaultCorsHeaders = map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Authorization, Content-Type, Content-Length",
			"Access-Control-Allow-Methods": "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD",
			"Access-Control-Expose-Headers": "Content-Length, Access-Control-Allow-Origin, " +
				"Access-Control-Allow-Headers, Content-Type, Authorization",
			"Access-Control-Allow-Credentials": "true",
		}
		if headers == nil {
			origin := c.Request.Header.Get("Origin")
			if len(origin) == 0 {
				origin = "*"
			}
			defaultCorsHeaders["Access-Control-Allow-Origin"] = origin
			headers = defaultCorsHeaders
		}
		for key, value := range headers {
			c.Writer.Header().Set(key, value)
		}

		c.Writer.WriteHeader(http.StatusNoContent)

		c.Next()
	}
}
