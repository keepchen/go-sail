package utils

import (
	"fmt"
	"strings"
)

var (
	summaryInfoTplDeprecated = `
// ----------- api doc definition -----------------

// @title          %s
// @version        %s
// @description   %s
// @termsOfService %s

// @contact.name  %s
// @contact.url   %s
// @contact.email %s

// @license.name MIT
// @license.url  %s

// @Host     %s
// @BasePath %s
%s
// ----------- api doc definition -----------------`
	needAuthorizeTplDeprecated = `// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
// @description                Access Token protects our entity endpoints`

	controllerInfoTplDeprecated = `
// %s %s
// @Tags        %s
// @Summary     %s
// @Description %s
%s
// @Accept      application/json
// @Produce     json
%s
// @Param %s %s true "parameters"
// @Success     200   {object} %s
// @Router      %s [%s]`

	controllerInfoSecurityTplDeprecated      = `// @Security    ApiKeyAuth`
	controllerInfoAuthorizationTplDeprecated = `// @Param Authorization header string true "Authorization token"`
	controllerInfoQueryTplDeprecated         = `query query`
	controllerInfoPayloadTplDeprecated       = `payload body`
)

// SwaggerSummaryInfoParamDeprecated 概览信息参数
type SwaggerSummaryInfoParamDeprecated struct {
	Title            string //标题
	Version          string //版本
	Description      string //描述
	TermOfServiceUrl string //服务条款地址
	ContactName      string //联系人
	ContactUrl       string //联系网址
	ContactEmail     string //联系邮箱
	LicenseUrl       string //证书网址
	Host             string //接口主机地址
	BasePath         string //接口公共前缀路径
	NeedAuthorize    bool   //是否需要token授权
}

// PrintSwaggerSummaryInfo 打印swagger概览信息参数
//
// Deprecated: PrintSwaggerSummaryInfo is deprecated,it will be removed in the future.
//
// Please use Swagger().PrintSummaryInfo() instead.
//
// 用于生成swagger注释
func PrintSwaggerSummaryInfo(param SwaggerSummaryInfoParamDeprecated) string {
	var needAuthorizeString string
	if param.NeedAuthorize {
		needAuthorizeString = needAuthorizeTplDeprecated
	}

	return fmt.Sprintf(summaryInfoTplDeprecated,
		param.Title, param.Version, param.Description, param.TermOfServiceUrl, param.ContactName, param.ContactUrl,
		param.ContactEmail, param.LicenseUrl, param.Host, param.BasePath, needAuthorizeString)
}

// SwaggerControllerInfoParamDeprecated 控制器信息参数
type SwaggerControllerInfoParamDeprecated struct {
	FunctionName       string //方法名称
	FunctionDesc       string //方法描述
	Tag                string //分类标签
	Summary            string //简要标题
	Description        string //接口描述
	Method             string //请求方法
	RequestParamString string //请求参数字符
	ResponseBodyString string //返回体字符
	ApiPath            string //接口路径
	NeedAuthorize      bool   //是否需要授权
}

// PrintSwaggerControllerInfo 打印swagger控制器信息
//
// Deprecated: PrintSwaggerControllerInfo is deprecated,it will be removed in the future.
//
// Please use Swagger().PrintControllerInfo() instead.
//
// 用于生成swagger注释
func PrintSwaggerControllerInfo(param SwaggerControllerInfoParamDeprecated) string {
	var (
		requestString       string
		securityString      string
		needAuthorizeString string
	)
	if param.NeedAuthorize {
		securityString = controllerInfoSecurityTplDeprecated
		needAuthorizeString = controllerInfoAuthorizationTplDeprecated
	}
	if len(param.Method) == 0 {
		param.Method = "get"
	}
	param.Method = strings.ToLower(param.Method)
	if param.Method != "get" {
		requestString = controllerInfoPayloadTplDeprecated
	} else {
		requestString = controllerInfoQueryTplDeprecated
	}

	return fmt.Sprintf(controllerInfoTplDeprecated,
		param.FunctionName, param.FunctionDesc, param.Tag, param.Summary, param.Description, securityString,
		needAuthorizeString, requestString, param.RequestParamString, param.ResponseBodyString, param.ApiPath, param.Method)
}
