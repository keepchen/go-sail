package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithCorsOnlyOptions 允许浏览器跨域请求，仅处理options探测请求
func WithCorsOnlyOptions(headers map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//处理浏览器跨域options探测请求
		if c.Request.Method == http.MethodOptions {
			origin := c.Request.Header.Get("Origin")
			if len(origin) == 0 {
				origin = "*"
			}

			if headers == nil {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length")
				c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD")
				c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
					"Access-Control-Allow-Headers, Content-Type, Authorization")
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			} else {
				for key, value := range headers {
					c.Writer.Header().Set(key, value)
				}
			}

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
		origin := c.Request.Header.Get("Origin")
		if len(origin) == 0 {
			origin = "*"
		}

		if headers == nil {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS, HEAD")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, "+
				"Access-Control-Allow-Headers, Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			for key, value := range headers {
				c.Writer.Header().Set(key, value)
			}
		}

		c.Writer.WriteHeader(http.StatusNoContent)

		c.Next()
	}
}
