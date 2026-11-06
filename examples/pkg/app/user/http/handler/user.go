package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/service"
)

// GetUserInfo retrieves user information
// @Tags        user / User operations
// @Summary     user-info / Get user information
// @Description Get user information
// @Security    ApiKeyAuth
// @Accept      application/json
// @Produce     json
// @Param       query query    request.GetUserInfo true "Query parameters"
// @Success     200   {object} response.GetUserInfo
// @Failure     400   {object} dto.Error400
// @Failure     500   {object} dto.Error500
// @Router      /user/info [get]
func GetUserInfo(c *gin.Context) {
	service.GetUserInfoSvc(c)
}
