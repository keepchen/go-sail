package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/keepchen/go-sail/pkg/app/user/http/vo/request"
	"github.com/keepchen/go-sail/pkg/app/user/http/vo/response"
	"github.com/keepchen/go-sail/pkg/common/http/api"
	"github.com/keepchen/go-sail/pkg/constants"
)

func SayHelloSvc(c *gin.Context) {
	var (
		form request.SayHello
		resp response.SayHello
	)
	if err := c.ShouldBind(&form); err != nil {
		api.New(c).Assemble(constants.ErrRequestParamsInvalid, nil).Send()
		return
	}

	if errorCode, err := form.Validator(); err != nil {
		api.New(c).Assemble(errorCode, nil, err.Error()).Send()
		return
	}

	var nickname string
	if len(form.Nickname) == 0 {
		nickname = "go-sail"
	} else {
		nickname = form.Nickname
	}

	resp.Data = fmt.Sprintf("hello, %s", nickname)

	api.New(c).Assemble(constants.ErrNone, resp).Send()
}
