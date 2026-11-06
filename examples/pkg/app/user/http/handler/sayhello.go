package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/v3/examples/pkg/app/user/service"
)

// SayHello greeting
// @Tags        welcome / Greeting
// @Summary     say-hello / Say hello
// @Description Greeting
// @Accept      application/json
// @Produce     json
// @Param       query query    request.SayHello true "Query parameters"
// @Success     200   {object} response.SayHello
// @Failure     400   {object} dto.Error400
// @Failure     500   {object} dto.Error500
// @Router      /say-hello [get]
func SayHello(c *gin.Context) {
	service.SayHelloSvc(c)
}
