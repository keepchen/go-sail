package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/pkg/app/user/service"
)

// GetUserInfo 获取用户信息
// @Tags        user / 用户相关
// @Summary     user-info / 获取用户信息
// @Description 获取用户信息
// @Security    ApiKeyAuth
// @Accept      application/json
// @Produce     json
// @Param       query query    request.GetUserInfo true "查询参数"
// @Success     200   {object} response.GetUserInfo
// @Failure     400   {object} api.Error400
// @Failure     500   {object} api.Error500
// @Router      /user/info [get]
func GetUserInfo(c *gin.Context) {
	service.GetUserInfoSvc(c)
}
