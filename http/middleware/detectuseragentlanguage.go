package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// DetectUserAgentLanguage 检测客户端语言
//
// 通过解析请求头中的`Accept-Language`字段得到
//
// 默认为`en-US`
func DetectUserAgentLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var language = "en-US"
		al := c.Request.Header.Get("accept-language")
		if len(al) > 0 {
			als := strings.Split(al, ",")
			if len(als) > 0 {
				language = als[0]
			}
		}
		c.Set("language", language)

		c.Next()
	}
}
