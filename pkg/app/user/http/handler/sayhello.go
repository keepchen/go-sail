package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v2/pkg/app/user/service"
)

// SayHello 欢迎语
// @Tags        welcome / 欢迎语
// @Summary     say-hello / 说你好
// @Description 欢迎语
// @Accept      application/json
// @Produce     json
// @Param       query query    request.SayHello true "查询参数"
// @Success     200   {object} response.SayHello
// @Failure     400   {object} dto.Error400
// @Failure     500   {object} dto.Error500
// @Router      /say-hello [get]
func SayHello(c *gin.Context) {
	service.SayHelloSvc(c)
}
