package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WithCorsOnlyOptions allows cross-origin browser requests
//
// Inject permissive headers only for OPTIONS preflight requests
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

		//Handle cross-origin browser OPTIONS preflight requests
		if c.Request.Method == http.MethodOptions {
			//Mobile browsers may disallow cross-origin requests returning 204 No Content, mainly for security and standards compliance.
			//The browser same-origin policy and CORS specification aim to ensure transparency and security for cross-origin requests,
			//while the lack of content in a 204 response may cause security or compatibility issues, so browsers may block it in cross-origin scenarios.
			//c.Writer.WriteHeader(http.StatusNoContent)
			c.Abort()
			return
		}

		c.Next()
	}
}

// WithCors allows cross-origin browser requests
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

		//Mobile browsers may disallow cross-origin requests returning 204 No Content, mainly for security and standards compliance.
		//The browser same-origin policy and CORS specification aim to ensure transparency and security for cross-origin requests,
		//while the lack of content in a 204 response may cause security or compatibility issues, so browsers may block it in cross-origin scenarios.
		//if c.Request.Method == http.MethodOptions {
		//	statusCode = http.StatusNoContent
		//}

		c.Writer.WriteHeader(statusCode)

		c.Next()
	}
}
