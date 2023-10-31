package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/constants"
	"github.com/keepchen/go-sail/v3/http/api"
)

// AuthCheck 授权检查
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if len(authorization) == 0 {
			api.New(c).Assemble(constants.ErrAuthorizationTokenInvalid, nil).SendWithCode(http.StatusUnauthorized)
			return
		}

		//more...

		c.Next()
	}
}
