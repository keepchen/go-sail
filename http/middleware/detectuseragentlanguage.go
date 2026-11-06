package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// DetectUserAgentLanguage detects the client language
//
// Obtained by parsing the `Accept-Language` request header
//
// Defaults to `en`
func DetectUserAgentLanguage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var language = "en"
		al := c.Request.Header.Get("accept-language")
		if len(al) > 0 {
			//example: zh-CN,zh;q=0.9,en;q=0.8,ja;q=0.7
			als := strings.Split(al, ",")
			if len(als) > 0 {
				language = als[0]
			}
		}
		c.Set("language", language)

		c.Next()
	}
}
