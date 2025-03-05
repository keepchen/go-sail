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
		var (
			statusCode         = http.StatusOK
			defaultCorsHeaders = map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Authorization, Content-Type, Content-Length",
				"Access-Control-Allow-Methods": "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD",
				"Access-Control-Expose-Headers": "Content-Length, Access-Control-Allow-Origin, " +
					"Access-Control-Allow-Headers, Content-Type, Authorization",
				"Access-Control-Allow-Credentials": "false",
			}
		)
		if headers == nil {
			origin := c.Request.Header.Get("Origin")
			if len(origin) == 0 {
				origin = "*"
			} else {
				defaultCorsHeaders["Access-Control-Allow-Credentials"] = "true"
			}
			defaultCorsHeaders["Access-Control-Allow-Origin"] = origin
			headers = defaultCorsHeaders
		}
		for key, value := range headers {
			c.Writer.Header().Set(key, value)
		}

		c.Writer.WriteHeader(statusCode)

		//处理浏览器跨域options探测请求
		if c.Request.Method == http.MethodOptions {
			//手机端浏览器中返回 204 No Content 状态的跨域请求不被允许，主要是出于安全性和规范性的考虑。
			//浏览器的同源策略和 CORS 规范旨在确保跨域请求的透明性和安全性，
			//而 204 响应的无内容特性可能导致安全隐患或兼容性问题，因此在跨域环境中可能会被浏览器禁止。
			//c.Writer.WriteHeader(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}

// WithCors 允许浏览器跨域请求
func WithCors(headers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			statusCode         = http.StatusOK
			defaultCorsHeaders = map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Headers": "Authorization, Content-Type, Content-Length",
				"Access-Control-Allow-Methods": "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD",
				"Access-Control-Expose-Headers": "Content-Length, Access-Control-Allow-Origin, " +
					"Access-Control-Allow-Headers, Content-Type, Authorization",
				"Access-Control-Allow-Credentials": "false",
			}
		)
		if headers == nil {
			origin := c.Request.Header.Get("Origin")
			if len(origin) == 0 {
				origin = "*"
			} else {
				defaultCorsHeaders["Access-Control-Allow-Credentials"] = "true"
			}
			defaultCorsHeaders["Access-Control-Allow-Origin"] = origin
			headers = defaultCorsHeaders
		}
		for key, value := range headers {
			c.Writer.Header().Set(key, value)
		}

		//手机端浏览器中返回 204 No Content 状态的跨域请求不被允许，主要是出于安全性和规范性的考虑。
		//浏览器的同源策略和 CORS 规范旨在确保跨域请求的透明性和安全性，
		//而 204 响应的无内容特性可能导致安全隐患或兼容性问题，因此在跨域环境中可能会被浏览器禁止。
		//if c.Request.Method == http.MethodOptions {
		//	statusCode = http.StatusNoContent
		//}

		c.Writer.WriteHeader(statusCode)

		c.Next()
	}
}
