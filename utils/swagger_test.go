package utils

import (
	"fmt"
	"testing"
)

func TestPrintSwaggerSummaryInfo(t *testing.T) {
	p0 := SwaggerSummaryInfoParam{}
	fmt.Println(PrintSwaggerSummaryInfo(p0))
	t.Log("OK")
	p1 := SwaggerSummaryInfoParam{
		Title:            "go-sail",
		Version:          "1.0",
		Description:      "go-sail是一个轻量级的web框架 / The go-sail is a lightweight golang web framework.",
		TermOfServiceUrl: "https://blog.keepchen.com",
		ContactName:      "keepchen",
		ContactUrl:       "https://blog.keepchen.com",
		ContactEmail:     "keepchen2016@gmail.com",
		LicenseUrl:       "https://github.com/keepchen/go-sail/LICENSE",
		Host:             "http://127.0.0.1",
		BasePath:         "/api/v1",
		NeedAuthorize:    true,
	}
	fmt.Println(PrintSwaggerSummaryInfo(p1))
	t.Log("OK")
}

func TestPrintSwaggerControllerInfo(t *testing.T) {
	p0 := SwaggerControllerInfoParam{
		FunctionName:       "GetUserInfo",
		FunctionDesc:       "获取用户信息",
		Tag:                "用户模块",
		Summary:            "获取用户信息",
		Description:        "根据授权信息获取用户的信息",
		Method:             "get",
		RequestParamString: "request.GetUserInfoVo",
		ResponseBodyString: "response.GetUserInfoDto",
		ApiPath:            "/api/v1/user/info/get",
		NeedAuthorize:      true,
	}
	fmt.Println(PrintSwaggerControllerInfo(p0))
	t.Log("OK")
	p1 := SwaggerControllerInfoParam{
		FunctionName:       "UpdateUserInfo",
		FunctionDesc:       "更新用户信息",
		Tag:                "用户模块",
		Summary:            "更新用户信息",
		Description:        "根据授权信息更新用户的信息",
		Method:             "post",
		RequestParamString: "request.UpdateUserInfoVo",
		ResponseBodyString: "response.UpdateUserInfoDto",
		ApiPath:            "/api/v1/user/info/update",
		NeedAuthorize:      true,
	}
	fmt.Println(PrintSwaggerControllerInfo(p1))
	t.Log("OK")
}
